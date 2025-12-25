// Package log provides QMUILog - a logging framework
// Ported from Tencent's QMUI_iOS framework
package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel defines the logging level
type LogLevel int

const (
	// LogLevelDefault is the default log level
	LogLevelDefault LogLevel = iota
	// LogLevelInfo is for informational messages
	LogLevelInfo
	// LogLevelWarn is for warning messages
	LogLevelWarn
	// LogLevelError is for error messages
	LogLevelError
)

// LogItem represents a single log entry
type LogItem struct {
	Level     LogLevel
	Name      string
	Message   string
	Timestamp time.Time
	File      string
	Line      int
}

// String returns a formatted string representation of the log item
func (li *LogItem) String() string {
	levelStr := ""
	switch li.Level {
	case LogLevelDefault:
		levelStr = "DEFAULT"
	case LogLevelInfo:
		levelStr = "INFO"
	case LogLevelWarn:
		levelStr = "WARN"
	case LogLevelError:
		levelStr = "ERROR"
	}

	return fmt.Sprintf("[%s] [%s] %s: %s",
		li.Timestamp.Format("2006-01-02 15:04:05.000"),
		levelStr,
		li.Name,
		li.Message,
	)
}

// Logger is a log handler
type Logger struct {
	Name    string
	Level   LogLevel
	Output  io.Writer
	Enabled bool

	mu       sync.RWMutex
	items    []*LogItem
	handlers []func(item *LogItem)
}

// NewLogger creates a new logger
func NewLogger(name string) *Logger {
	return &Logger{
		Name:     name,
		Level:    LogLevelDefault,
		Output:   os.Stdout,
		Enabled:  true,
		items:    make([]*LogItem, 0),
		handlers: make([]func(item *LogItem), 0),
	}
}

// Log logs a message at the specified level
func (l *Logger) Log(level LogLevel, message string) {
	if !l.Enabled || level < l.Level {
		return
	}

	item := &LogItem{
		Level:     level,
		Name:      l.Name,
		Message:   message,
		Timestamp: time.Now(),
	}

	l.mu.Lock()
	l.items = append(l.items, item)
	handlers := l.handlers
	output := l.Output
	l.mu.Unlock()

	// Write to output
	if output != nil {
		fmt.Fprintln(output, item.String())
	}

	// Call handlers
	for _, handler := range handlers {
		handler(item)
	}
}

// LogDefault logs at default level
func (l *Logger) LogDefault(format string, args ...interface{}) {
	l.Log(LogLevelDefault, fmt.Sprintf(format, args...))
}

// LogInfo logs at info level
func (l *Logger) LogInfo(format string, args ...interface{}) {
	l.Log(LogLevelInfo, fmt.Sprintf(format, args...))
}

// LogWarn logs at warn level
func (l *Logger) LogWarn(format string, args ...interface{}) {
	l.Log(LogLevelWarn, fmt.Sprintf(format, args...))
}

// LogError logs at error level
func (l *Logger) LogError(format string, args ...interface{}) {
	l.Log(LogLevelError, fmt.Sprintf(format, args...))
}

// AddHandler adds a log handler
func (l *Logger) AddHandler(handler func(item *LogItem)) {
	l.mu.Lock()
	l.handlers = append(l.handlers, handler)
	l.mu.Unlock()
}

// GetItems returns all logged items
func (l *Logger) GetItems() []*LogItem {
	l.mu.RLock()
	defer l.mu.RUnlock()
	items := make([]*LogItem, len(l.items))
	copy(items, l.items)
	return items
}

// ClearItems clears all logged items
func (l *Logger) ClearItems() {
	l.mu.Lock()
	l.items = make([]*LogItem, 0)
	l.mu.Unlock()
}

// LogNameManager manages multiple loggers by name
type LogNameManager struct {
	loggers map[string]*Logger
	mu      sync.RWMutex
}

// NewLogNameManager creates a new log name manager
func NewLogNameManager() *LogNameManager {
	return &LogNameManager{
		loggers: make(map[string]*Logger),
	}
}

// GetLogger gets or creates a logger by name
func (m *LogNameManager) GetLogger(name string) *Logger {
	m.mu.RLock()
	logger, exists := m.loggers[name]
	m.mu.RUnlock()

	if exists {
		return logger
	}

	m.mu.Lock()
	logger = NewLogger(name)
	m.loggers[name] = logger
	m.mu.Unlock()

	return logger
}

// GetAllLoggers returns all loggers
func (m *LogNameManager) GetAllLoggers() []*Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()

	loggers := make([]*Logger, 0, len(m.loggers))
	for _, l := range m.loggers {
		loggers = append(loggers, l)
	}
	return loggers
}

// Global logger manager
var (
	sharedLogManager *LogNameManager
	logManagerOnce   sync.Once
)

// SharedLogManager returns the shared log manager
func SharedLogManager() *LogNameManager {
	logManagerOnce.Do(func() {
		sharedLogManager = NewLogNameManager()
	})
	return sharedLogManager
}

// Convenience functions using a default "QMUI" logger

// QMUILog logs at default level
func QMUILog(format string, args ...interface{}) {
	SharedLogManager().GetLogger("QMUI").LogDefault(format, args...)
}

// QMUILogInfo logs at info level
func QMUILogInfo(format string, args ...interface{}) {
	SharedLogManager().GetLogger("QMUI").LogInfo(format, args...)
}

// QMUILogWarn logs at warn level
func QMUILogWarn(format string, args ...interface{}) {
	SharedLogManager().GetLogger("QMUI").LogWarn(format, args...)
}

// QMUILogError logs at error level
func QMUILogError(format string, args ...interface{}) {
	SharedLogManager().GetLogger("QMUI").LogError(format, args...)
}

// SetEnabled enables/disables the QMUI logger
func SetEnabled(enabled bool) {
	SharedLogManager().GetLogger("QMUI").Enabled = enabled
}

// SetLevel sets the minimum log level
func SetLevel(level LogLevel) {
	SharedLogManager().GetLogger("QMUI").Level = level
}

// SetOutput sets the output writer
func SetOutput(output io.Writer) {
	SharedLogManager().GetLogger("QMUI").Output = output
}

// FileLogger writes logs to a file
type FileLogger struct {
	*Logger
	FilePath string
	file     *os.File
	mu       sync.Mutex
}

// NewFileLogger creates a file logger
func NewFileLogger(name, filePath string) (*FileLogger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	fl := &FileLogger{
		Logger:   NewLogger(name),
		FilePath: filePath,
		file:     file,
	}
	fl.Output = file

	return fl, nil
}

// Close closes the file
func (fl *FileLogger) Close() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	if fl.file != nil {
		return fl.file.Close()
	}
	return nil
}

// RotateFile rotates the log file
func (fl *FileLogger) RotateFile() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	if fl.file != nil {
		fl.file.Close()
	}

	// Rename current file
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fl.FilePath + "." + timestamp
	os.Rename(fl.FilePath, backupPath)

	// Create new file
	file, err := os.OpenFile(fl.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	fl.file = file
	fl.Output = file
	return nil
}
