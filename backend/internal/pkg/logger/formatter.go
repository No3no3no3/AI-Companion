package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// CustomFormatter 自定义格式化器，实现logrus.Formatter接口
type CustomFormatter struct {
	// 是否显示颜色
	EnableColor bool
	// 时间戳格式
	TimestampFormat string
	// 是否显示完整时间戳
	FullTimestamp bool
	// 颜色映射
	Colors map[logrus.Level]string
}

// NewCustomFormatter 创建自定义格式化器
func NewCustomFormatter(enableColor bool) *CustomFormatter {
	return &CustomFormatter{
		EnableColor:     enableColor,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		Colors: map[logrus.Level]string{
			logrus.DebugLevel: "\033[36m", // 青色
			logrus.InfoLevel:  "\033[32m", // 绿色
			logrus.WarnLevel:  "\033[33m", // 黄色
			logrus.ErrorLevel: "\033[31m", // 红色
			logrus.FatalLevel: "\033[31m", // 红色
			logrus.PanicLevel: "\033[31m", // 红色
		},
	}
}

// Format 实现logrus.Formatter接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 添加时间戳
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if f.EnableColor {
		// 时间为绿色
		b.WriteString("\033[32m")
	}
	b.WriteString(entry.Time.Format(timestampFormat))
	if f.EnableColor {
		b.WriteString("\033[0m")
	}
	b.WriteString(" ")

	// 添加日志等级
	levelText := strings.ToUpper(entry.Level.String())
	if f.EnableColor {
		if color, ok := f.Colors[entry.Level]; ok {
			b.WriteString(color)
		}
	}
	b.WriteString("[")
	b.WriteString(levelText)
	b.WriteString("]")
	if f.EnableColor {
		b.WriteString("\033[0m")
	}
	b.WriteString(" ")

	// 添加函数和行号信息（如果有）
	if entry.HasCaller() {
		fname := filepathBase(entry.Caller.Function)
		if f.EnableColor {
			// 函数名为青色
			b.WriteString("\033[36m")
		}
		b.WriteString("[")
		b.WriteString(fname)
		b.WriteString(":")
		b.WriteString(fmt.Sprintf("%d", entry.Caller.Line))
		b.WriteString("]")
		if f.EnableColor {
			b.WriteString("\033[0m")
		}
		b.WriteString(" ")
	}

	// 添加消息
	b.WriteString(entry.Message)

	// 添加字段
	for _, key := range keys {
		b.WriteString(" ")
		b.WriteString(key)
		b.WriteString("=")
		f.appendValue(b, entry.Data[key])
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// appendValue 追加值到缓冲区
func (f *CustomFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

// needsQuoting 检查字符串是否需要引号
func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

// filepathBase 从完整路径中提取文件名
func filepathBase(path string) string {
	if path == "" {
		return ""
	}

	// 分割路径获取最后一个元素
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		// 提取函数名
		funcPart := parts[len(parts)-1]
		// 处理包名和函数名
		if lastDot := strings.LastIndex(funcPart, "."); lastDot >= 0 {
			return funcPart[lastDot+1:]
		}
		return funcPart
	}
	return path
}

// 默认时间戳格式
const defaultTimestampFormat = "2006-01-02 15:04:05"

// CallerHook 调用者信息Hook
type CallerHook struct {
	Skip int // 跳过的调用栈层数
}

// Levels 返回该Hook处理的日志等级
func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 实现logrus.Hook接口
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	if entry.HasCaller() {
		// 已经有调用者信息，不需要额外处理
		return nil
	}

	// 获取调用者信息
	pc, file, line, ok := runtime.Caller(hook.Skip)
	if !ok {
		return nil
	}

	// 设置调用者信息
	entry.Caller = &runtime.Frame{
		PC:   pc,
		File: file,
		Line: line,
		Func: runtime.FuncForPC(pc),
	}

	return nil
}

// JSONFormatterWithCaller 带调用者信息的JSON格式化器
type JSONFormatterWithCaller struct {
	logrus.JSONFormatter
}

// NewJSONFormatterWithCaller 创建带调用者信息的JSON格式化器
func NewJSONFormatterWithCaller() *JSONFormatterWithCaller {
	return &JSONFormatterWithCaller{
		JSONFormatter: logrus.JSONFormatter{
			TimestampFormat: defaultTimestampFormat,
		},
	}
}

// Format 实现logrus.Formatter接口
func (f *JSONFormatterWithCaller) Format(entry *logrus.Entry) ([]byte, error) {
	// 确保调用者信息被包含
	if !entry.HasCaller() {
		return f.JSONFormatter.Format(entry)
	}

	// 复制数据以避免修改原始entry
	data := make(logrus.Fields, len(entry.Data)+4)
	for k, v := range entry.Data {
		data[k] = v
	}

	// 添加调用者信息
	if entry.Caller != nil {
		if entry.Caller.Func != nil {
			data["function"] = entry.Caller.Func.Name()
		}
		data["file"] = entry.Caller.File
		data["line"] = entry.Caller.Line
	}

	// 创建新entry
	newEntry := &logrus.Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
	}

	return f.JSONFormatter.Format(newEntry)
}
