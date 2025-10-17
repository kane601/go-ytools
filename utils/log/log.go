package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Level represents log severity.
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return fmt.Sprintf("LEVEL(%d)", int(l))
	}
}

// Logger is a simple leveled logger with a standard library backend.
type Logger struct {
	mu    sync.Mutex
	level Level
	std   *log.Logger
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func initDefault() {
	defaultLogger = &Logger{
		level: InfoLevel,
		std:   log.New(os.Stderr, "", 0),
	}
}

// SetLevel sets the global log level.
func SetLevel(level Level) { get().SetLevel(level) }

// GetLevel gets the global log level.
func GetLevel() Level { return get().GetLevel() }

// WithOutput replaces the output writer for the global logger.
func WithOutput(file *os.File) { get().WithOutput(file) }

// get returns the singleton logger instance.
func get() *Logger {
	once.Do(initDefault)
	return defaultLogger
}

// New returns a new Logger with the given level and output.
func New(level Level, out *os.File) *Logger {
	if out == nil {
		out = os.Stderr
	}
	return &Logger{level: level, std: log.New(out, "", 0)}
}

// SetLevel sets the logger level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns current logger level.
func (l *Logger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// WithOutput changes the output writer.
func (l *Logger) WithOutput(file *os.File) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.std.SetOutput(file)
}

func (l *Logger) logf(level Level, format string, args ...any) {
	l.mu.Lock()
	current := l.level
	l.mu.Unlock()
	if level < current {
		return
	}
	ts := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf(format, args...)
	l.std.Printf("%s [%s] %s", ts, level.String(), msg)
}

// Global helpers
func Debugf(format string, args ...any) { get().logf(DebugLevel, format, args...) }
func Infof(format string, args ...any)  { get().logf(InfoLevel, format, args...) }
func Warnf(format string, args ...any)  { get().logf(WarnLevel, format, args...) }
func Errorf(format string, args ...any) { get().logf(ErrorLevel, format, args...) }

// Non-format variants
func Debug(args ...any) { Debugf("%v", fmt.Sprint(args...)) }
func Info(args ...any)  { Infof("%v", fmt.Sprint(args...)) }
func Warn(args ...any)  { Warnf("%v", fmt.Sprint(args...)) }
func Error(args ...any) { Errorf("%v", fmt.Sprint(args...)) }

