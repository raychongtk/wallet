package util

import (
	"go.uber.org/zap"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitializeLogger initializes the logger as a singleton
func InitializeLogger(isProduction bool) {
	once.Do(func() {
		var err error
		if isProduction {
			logger, err = zap.NewProduction()
		} else {
			logger, err = zap.NewDevelopment()
		}
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
	})
}

// GetLogger returns the singleton logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		panic("logger is not initialized. Call InitializeLogger first.")
	}
	return logger
}

// Info logs an info-level message
func Info(message string, fields ...zap.Field) {
	GetLogger().Info(message, fields...)
}

// Error logs an error-level message
func Error(message string, fields ...zap.Field) {
	GetLogger().Error(message, fields...)
}

// Debug logs a debug-level message
func Debug(message string, fields ...zap.Field) {
	GetLogger().Debug(message, fields...)
}

// Warn logs a warning-level message
func Warn(message string, fields ...zap.Field) {
	GetLogger().Warn(message, fields...)
}
