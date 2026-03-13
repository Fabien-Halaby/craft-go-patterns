package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton/log"
)

func main() {
	log := logger.GetInstance(logger.Options{
		Level:     logger.DebugLevel,
		Formatter: &logger.JSONFormatter{},
		Output:    os.Stdout,
		Caller:    true,
		Fields:    logger.Fields{logger.String("app", "demo"), logger.String("version", "1.0.0")},
	})

	log.Info("Server starting", logger.Int("port", 8080))

	reqLog := log.With(logger.String("request_id", "abc-123"), logger.String("user", "fabien"))
	
	reqLog.Info("Handling request", logger.String("path", "/api/users"))
	
	start := time.Now()

	reqLog.Info("Request completed", 
		logger.Duration("duration", time.Since(start)),
		logger.Int("status", 200),
	)

	if err := doSomething(); err != nil {
		reqLog.Error("Operation failed", logger.Error(err))
	}

	devLog := logger.GetInstance(logger.Options{
		Level:     logger.DebugLevel,
		Formatter: &logger.TextFormatter{},
		Output:    os.Stdout,
	})
	devLog.Debug("Debug mode", logger.Bool("verbose", true))
}

func doSomething() error {
	return fmt.Errorf("connection refused")
}