package logger

import (
	"time"
)

/**
	Field:
		- represents a structured log field
**/
type Field struct {
	Key   string
	Value any
}

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field { 
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

func Duration(key string, d time.Duration) Field {
	return Field{Key: key, Value: d.String()}
}

type Fields []Field

func (f Fields) Map() map[string]any {
	m := make(map[string]any, len(f))
	for _, field := range f {
		m[field.Key] = field.Value
	}
	return m
}

/**
	With:
		- creates new Fields with additional fields
**/
func (f Fields) With(fields ...Field) Fields {
	return append(append(Fields(nil), f...), fields...)
}