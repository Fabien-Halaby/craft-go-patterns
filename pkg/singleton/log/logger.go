package logger

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton"
)

/**
	Logger is the structured logging interface
**/
type Logger struct {
	level     atomic.Int32
	formatter Formatter
	output    io.Writer
	mu        sync.Mutex
	fields    Fields
	caller    bool
}

type Options struct {
	Level     Level
	Formatter Formatter
	Output    io.Writer
	Fields    Fields
	Caller    bool
}

var instances = &sync.Map{}

/**
	GetInstance:
		- returns singleton Logger with options
**/
func GetInstance(opts Options) *Logger {
	key := reflect.TypeOf((*Logger)(nil))

	actual, _ := instances.LoadOrStore(key, singleton.NewInstance(func() *Logger {
		if opts.Formatter == nil {
			opts.Formatter = &TextFormatter{}
		}
		if opts.Output == nil {
			opts.Output = os.Stdout
		}

		l := &Logger{
			formatter: opts.Formatter,
			output:    opts.Output,
			fields:    opts.Fields,
			caller:    opts.Caller,
		}
		l.level.Store(int32(opts.Level))
		return l
	}))

	return actual.(*singleton.Instance[*Logger]).Get()
}

/**
	MustGetInstance panics if setup fails
**/
func MustGetInstance(opts Options) *Logger {
	return GetInstance(opts)
}

/**
	Default returns logger with sensible defaults
**/
func Default() *Logger {
	return GetInstance(Options{
		Level:     InfoLevel,
		Formatter: &TextFormatter{},
		Output:    os.Stdout,
	})
}

/**
	With creates child logger with additional fields
**/
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{
		level:     atomic.Int32{},
		formatter: l.formatter,
		output:    l.output,
		fields:    l.fields.With(fields...),
		caller:    l.caller,
	}
}

/**
	SetLevel changes log level dynamically
**/
func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

/**
	Core logging method
**/
func (l *Logger) log(level Level, msg string, fields Fields) {
	if level < Level(l.level.Load()) {
		return
	}

	entry := Entry{
		Time:    time.Now().UTC(),
		Level:   level,
		Message: msg,
		Fields:  l.fields.With(fields...),
	}

	if l.caller {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			parts := strings.Split(file, "/")
			if len(parts) > 2 {
				file = strings.Join(parts[len(parts)-2:], "/")
			}
			entry.Caller = fmt.Sprintf("%s:%d", file, line)
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.formatter.Format(entry, l.output)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields)
	os.Exit(1)
}

/**
	Sugar methods (printf-style)
**/
func (l *Logger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}
func (l *Logger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

/**
	ResetForTest clears singleton
**/
func ResetForTest() {
	key := reflect.TypeOf((*Logger)(nil))
	instances.Delete(key)
}