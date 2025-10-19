package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config logger配置结构
type Config struct {
	Level       LogLevel `json:"level"`        // 日志等级
	OutputDir   string   `json:"output_dir"`   // 日志输出目录
	Console     bool     `json:"console"`      // 是否输出到控制台
	File        bool     `json:"file"`         // 是否输出到文件
	MaxSize     int      `json:"max_size"`     // 单个日志文件最大大小(MB)
	MaxBackups  int      `json:"max_backups"`  // 保留的最大日志文件数
	MaxAge      int      `json:"max_age"`      // 日志文件保留天数
	EnableColor bool     `json:"enable_color"` // 是否启用颜色
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Level:       INFO,
		OutputDir:   "logs",
		Console:     true,
		File:        true,
		MaxSize:     100, // 100MB
		MaxBackups:  30,  // 保留30个文件
		MaxAge:      7,   // 保留7天
		EnableColor: true,
	}
}

// LogLevel 自定义日志等级类型
type LogLevel logrus.Level

const (
	// DEBUG 调试级别
	DEBUG LogLevel = LogLevel(logrus.DebugLevel)
	// INFO 信息级别
	INFO LogLevel = LogLevel(logrus.InfoLevel)
	// WARN 警告级别
	WARN LogLevel = LogLevel(logrus.WarnLevel)
	// ERROR 错误级别
	ERROR LogLevel = LogLevel(logrus.ErrorLevel)
)

// String 返回日志等级的字符串表示
func (l LogLevel) String() string {
	return logrus.Level(l).String()
}

// Logger 基于 logrus 的日志结构
type Logger struct {
	logger     *logrus.Logger
	entry      *logrus.Entry
	config     Config
	callerHook *CallerHook
}

// New 创建新的Logger实例
func New(config Config) (*Logger, error) {
	logger := logrus.New()

	// 设置日志等级
	logger.SetLevel(logrus.Level(config.Level))

	// 设置报告调用者信息
	logger.SetReportCaller(true)

	// 创建自定义调用者Hook
	callerHook := &CallerHook{
		Skip: 5, // 跳过堆栈层数，以获取正确的调用者信息
	}
	logger.AddHook(callerHook)

	// 创建自定义格式化器
	formatter := &CustomFormatter{
		EnableColor:     config.EnableColor,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	// 设置输出
	var writers []interface{}

	// 控制台输出
	if config.Console {
		writers = append(writers, os.Stdout)
	}

	// 文件输出
	if config.File {
		// 确保输出目录存在
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// 生成按天的日志文件名
		fileName := fmt.Sprintf("log-%s.log", time.Now().Format("2006-0102"))
		filePath := filepath.Join(config.OutputDir, fileName)

		// 使用 lumberjack 进行文件轮转
		fileWriter := &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   true,
			LocalTime:  true,
		}

		writers = append(writers, fileWriter)
	}

	// 设置多路输出
	if len(writers) > 1 {
		// 如果需要同时输出到多个地方，需要实现自定义的多路写入器
		// 这里简化处理，优先使用控制台输出
		logger.SetOutput(writers[0].(interface {
			Write([]byte) (int, error)
		}))
	} else if len(writers) == 1 {
		logger.SetOutput(writers[0].(interface {
			Write([]byte) (int, error)
		}))
	}

	logInstance := &Logger{
		logger:     logger,
		config:     config,
		callerHook: callerHook,
	}

	// 创建基础entry，可以添加通用字段
	logInstance.entry = logger.WithFields(logrus.Fields{})

	return logInstance, nil
}

// NewDefault 创建使用默认配置的Logger
func NewDefault() (*Logger, error) {
	return New(DefaultConfig())
}

// SetLevel 设置日志等级
func (l *Logger) SetLevel(level LogLevel) {
	l.config.Level = level
	l.logger.SetLevel(logrus.Level(level))
}

// GetLevel 获取当前日志等级
func (l *Logger) GetLevel() LogLevel {
	return LogLevel(l.logger.GetLevel())
}

// WithFields 添加字段
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{
		logger:     l.logger,
		entry:      l.entry.WithFields(fields),
		config:     l.config,
		callerHook: l.callerHook,
	}
}

// WithField 添加单个字段
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger:     l.logger,
		entry:      l.entry.WithField(key, value),
		config:     l.config,
		callerHook: l.callerHook,
	}
}

// Debug 记录调试级别日志
func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugf 记录调试级别格式化日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info 记录信息级别日志
func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infof 记录信息级别格式化日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnf 记录警告级别格式化日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error 记录错误级别日志
func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorf 记录错误级别格式化日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal 记录致命错误级别日志并退出程序
func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Fatalf 记录致命错误级别格式化日志并退出程序
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// Panic 记录panic级别日志并panic
func (l *Logger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

// Panicf 记录panic级别格式化日志并panic
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

// GetConfig 获取当前配置
func (l *Logger) GetConfig() Config {
	return l.config
}

// Close 关闭Logger
func (l *Logger) Close() error {
	// logrus 和 lumberjack 会自动处理关闭
	return nil
}

// 全局Logger实例
var defaultLogger *Logger
var defaultLoggerOnce sync.Once

// InitDefaultLogger 初始化默认Logger
func InitDefaultLogger(config Config) error {
	var err error
	defaultLoggerOnce.Do(func() {
		defaultLogger, err = New(config)
	})
	return err
}

// GetDefaultLogger 获取默认Logger实例
func GetDefaultLogger() *Logger {
	if defaultLogger == nil {
		// 如果没有初始化，使用默认配置创建一个
		defaultLogger, _ = NewDefault()
	}
	return defaultLogger
}

// 全局便捷函数
func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

//func Infof(format string, args ...interface{}) {
//	GetDefaultLogger().Infof(format, args...)
//}

func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}
func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(args...)
}

func Panic(args ...interface{}) {
	GetDefaultLogger().Panic(args...)
}

func WithFields(fields map[string]interface{}) *Logger {
	return GetDefaultLogger().WithFields(fields)
}

func WithField(key string, value interface{}) *Logger {
	return GetDefaultLogger().WithField(key, value)
}
