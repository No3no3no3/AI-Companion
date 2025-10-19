package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DailyRotateHook 按天轮转的Hook
type DailyRotateHook struct {
	outputDir   string
	maxSize     int
	maxBackups  int
	maxAge      int
	currentDate string
	writer      io.Writer
	mu          sync.RWMutex
}

// NewDailyRotateHook 创建按天轮转的Hook
func NewDailyRotateHook(outputDir string, maxSize, maxBackups, maxAge int) *DailyRotateHook {
	return &DailyRotateHook{
		outputDir:  outputDir,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
	}
}

// Levels 返回该Hook处理的日志等级
func (hook *DailyRotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 实现logrus.Hook接口
func (hook *DailyRotateHook) Fire(entry *logrus.Entry) error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	// 检查是否需要轮转文件
	currentDate := entry.Time.Format("2006-01-02")
	if hook.currentDate != currentDate || hook.writer == nil {
		if err := hook.rotateFile(currentDate); err != nil {
			return fmt.Errorf("failed to rotate log file: %w", err)
		}
		hook.currentDate = currentDate
	}

	// 格式化日志条目
	formatter := &CustomFormatter{
		EnableColor:     false, // 文件输出不使用颜色
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}

	bytes, err := formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("failed to format log entry: %w", err)
	}

	// 写入文件
	if hook.writer != nil {
		_, err = hook.writer.Write(bytes)
		return err
	}

	return fmt.Errorf("no writer available")
}

// rotateFile 轮转文件
func (hook *DailyRotateHook) rotateFile(date string) error {
	// 关闭当前的writer
	if closer, ok := hook.writer.(io.Closer); ok {
		closer.Close()
	}

	// 确保输出目录存在
	if err := os.MkdirAll(hook.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 生成新的文件名
	fileName := fmt.Sprintf("log-%s.log", date)
	filePath := filepath.Join(hook.outputDir, fileName)

	// 创建新的lumberjack writer
	hook.writer = &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    hook.maxSize,
		MaxBackups: hook.maxBackups,
		MaxAge:     hook.maxAge,
		Compress:   true,
		LocalTime:  true,
	}

	return nil
}

// Close 关闭Hook
func (hook *DailyRotateHook) Close() error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	if closer, ok := hook.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// MultiWriterHook 多路输出Hook
type MultiWriterHook struct {
	writers []logrus.Hook
}

// NewMultiWriterHook 创建多路输出Hook
func NewMultiWriterHook(writers ...logrus.Hook) *MultiWriterHook {
	return &MultiWriterHook{
		writers: writers,
	}
}

// Levels 返回该Hook处理的日志等级
func (hook *MultiWriterHook) Levels() []logrus.Level {
	// 返回所有writer支持的等级的并集
	levelSet := make(map[logrus.Level]bool)
	for _, writer := range hook.writers {
		for _, level := range writer.Levels() {
			levelSet[level] = true
		}
	}

	var levels []logrus.Level
	for level := range levelSet {
		levels = append(levels, level)
	}

	return levels
}

// Fire 实现logrus.Hook接口
func (hook *MultiWriterHook) Fire(entry *logrus.Entry) error {
	var firstErr error
	for _, writer := range hook.writers {
		if err := writer.Fire(entry); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// ConsoleHook 控制台输出Hook
type ConsoleHook struct {
	formatter logrus.Formatter
}

// NewConsoleHook 创建控制台输出Hook
func NewConsoleHook(enableColor bool) *ConsoleHook {
	return &ConsoleHook{
		formatter: NewCustomFormatter(enableColor),
	}
}

// Levels 返回该Hook处理的日志等级
func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 实现logrus.Hook接口
func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	bytes, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}

	// 直接写入标准输出
	os.Stdout.Write(bytes)
	return nil
}

// AsyncHook 异步Hook
type AsyncHook struct {
	entryChan chan *logrus.Entry
	hook      logrus.Hook
	done      chan struct{}
	wg        sync.WaitGroup
}

// NewAsyncHook 创建异步Hook
func NewAsyncHook(hook logrus.Hook, bufferSize int) *AsyncHook {
	asyncHook := &AsyncHook{
		entryChan: make(chan *logrus.Entry, bufferSize),
		hook:      hook,
		done:      make(chan struct{}),
	}

	// 启动goroutine处理日志条目
	asyncHook.wg.Add(1)
	go asyncHook.processEntries()

	return asyncHook
}

// Levels 返回该Hook处理的日志等级
func (hook *AsyncHook) Levels() []logrus.Level {
	return hook.hook.Levels()
}

// Fire 实现logrus.Hook接口
func (hook *AsyncHook) Fire(entry *logrus.Entry) error {
	select {
	case hook.entryChan <- entry:
		return nil
	case <-hook.done:
		return fmt.Errorf("hook is closed")
	default:
		// 缓冲区满，直接写入
		return hook.hook.Fire(entry)
	}
}

// processEntries 处理日志条目
func (hook *AsyncHook) processEntries() {
	defer hook.wg.Done()

	for {
		select {
		case entry := <-hook.entryChan:
			hook.hook.Fire(entry)
		case <-hook.done:
			// 处理剩余的条目
			for len(hook.entryChan) > 0 {
				entry := <-hook.entryChan
				hook.hook.Fire(entry)
			}
			return
		}
	}
}

// Close 关闭异步Hook
func (hook *AsyncHook) Close() error {
	close(hook.done)
	hook.wg.Wait()

	// 如果底层hook支持关闭，则关闭它
	if closer, ok := hook.hook.(interface{ Close() error }); ok {
		return closer.Close()
	}

	return nil
}
