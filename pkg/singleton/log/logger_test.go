package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevels(t *testing.T) {
	ResetForTest()
	
	var buf bytes.Buffer
	log := GetInstance(Options{
		Level:     WarnLevel,
		Formatter: &JSONFormatter{},
		Output:    &buf,
	})

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	output := buf.String()
	assert.NotContains(t, output, "debug")
	assert.NotContains(t, output, "info")
	assert.Contains(t, output, "warn")
	assert.Contains(t, output, "error")
}

func TestWithFields(t *testing.T) {
	ResetForTest()
	
	var buf bytes.Buffer
	log := GetInstance(Options{
		Level:     InfoLevel,
		Formatter: &JSONFormatter{},
		Output:    &buf,
		Fields:    Fields{String("service", "test")},
	})

	child := log.With(String("trace_id", "xyz"))
	child.Info("message", Int("count", 42))

	output := buf.String()
	assert.Contains(t, output, `"service":"test"`)
	assert.Contains(t, output, `"trace_id":"xyz"`)
	assert.Contains(t, output, `"count":42`)
}

func TestSingleton(t *testing.T) {
	ResetForTest()
	
	log1 := Default()
	log2 := Default()
	
	assert.Same(t, log1, log2)
}

func TestCaller(t *testing.T) {
	ResetForTest()
	
	var buf bytes.Buffer
	log := GetInstance(Options{
		Level:     InfoLevel,
		Formatter: &JSONFormatter{},
		Output:    &buf,
		Caller:    true,
	})

	log.Info("test")
	
	output := buf.String()
	assert.Contains(t, output, `"caller":`)
}

func TestSugarMethods(t *testing.T) {
	ResetForTest()
	
	var buf bytes.Buffer
	log := GetInstance(Options{
		Level:     InfoLevel,
		Formatter: &TextFormatter{},
		Output:    &buf,
	})

	log.Infof("Hello %s, count=%d", "world", 42)
	
	output := buf.String()
	assert.Contains(t, output, "Hello world, count=42")
}