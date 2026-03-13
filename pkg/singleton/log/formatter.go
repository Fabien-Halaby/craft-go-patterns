package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

/**
	Formatter:
		- defines log output format
**/
type Formatter interface {
	Format(entry Entry, w io.Writer) error
}

/**
	Entry:
		- represents a log record
**/
type Entry struct {
	Time    time.Time
	Level   Level
	Message string
	Fields  Fields
	Caller  string
}

/**
	JSONFormatter:
		- outputs structured JSON
**/
type JSONFormatter struct{}

func (f *JSONFormatter) Format(e Entry, w io.Writer) error {
	m := map[string]any{
		"time":    e.Time.Format(time.RFC3339Nano),
		"level":   e.Level,
		"message": e.Message,
	}
	if e.Caller != "" {
		m["caller"] = e.Caller
	}
	for k, v := range e.Fields.Map() {
		m[k] = v
	}
	return json.NewEncoder(w).Encode(m)
}

/**
	TextFormatter:
		- outputs human-readable logs
**/
type TextFormatter struct{}

func (f *TextFormatter) Format(e Entry, w io.Writer) error {
	fields := ""
	for k, v := range e.Fields.Map() {
		fields += fmt.Sprintf(" %s=%v", k, v)
	}
	caller := ""
	if e.Caller != "" {
		caller = fmt.Sprintf(" [%s]", e.Caller)
	}
	_, err := fmt.Fprintf(w, "%s | %-5s | %s%s%s\n",
		e.Time.Format("15:04:05.000"),
		e.Level,
		e.Message,
		caller,
		fields,
	)
	return err
}