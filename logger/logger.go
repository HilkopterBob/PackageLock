package logger

import (
	"os"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// InitLogger initializes the zap.Logger once and returns the instance.
func InitLogger() (*zap.SugaredLogger, error) {
	err := os.MkdirAll("logs/", os.ModePerm)
	if err != nil {
		pp.Printf("Couldn't create 'logs' directory. Got: %s", err)
		panic(err)
	}

	// Ensure the logger is initialized only once
	loggerConfig := zap.Config{
		Encoding:         "console", // You can also use "json"
		OutputPaths:      []string{"logs/app.log"},
		ErrorOutputPaths: []string{"logs/app_error.log"},
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	// Build the logger
	logger, err := loggerConfig.Build()
	if err != nil {
		// Handle logger initialization error (no panic)
		logger = nil
	}
	// Return the initialized logger and any error that occurred
	return logger.Sugar(), err
}
