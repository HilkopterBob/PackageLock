package logger

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Module exports the logger module.
var Module = fx.Provide(NewLogger)

// NewLogger constructs a new logger instance.
func NewLogger() (*zap.Logger, error) {
	createLogDir()

	// Initialize Lumberjack logger for log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log", // Log file name
		MaxSize:    250,            // Max size in MB before rotation
		MaxBackups: 3,              // Max number of old log files to retain
		MaxAge:     28,             // Max number of days to retain old log files
		Compress:   true,           // Compress the rotated log files
	}


	// Encoder configuration
	encoderConfig := zapcore.EncoderConfig{
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

	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig) // Use NewJSONEncoder for JSON logs

	// Create WriteSyncer for Lumberjack logger
	fileWriter := zapcore.AddSync(lumberjackLogger)
	consoleWriter := zapcore.AddSync(os.Stdout)
	writeSyncer := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)

	logLevel := zapcore.InfoLevel

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core,
		zap.AddCaller(),                       // Include caller information
		zap.AddStacktrace(zapcore.ErrorLevel), // Include stacktrace for error-level logs
	)

	logger.Info("--------------------------------------------")
	logger.Info("Logger initialized successfully")

	go func(logger *zap.Logger) {
		for {
			logger.Info("Loop")
		}
	}(logger)

	return logger, nil
}

// Creates to logs/ directory
func createLogDir() {
	err := os.MkdirAll("logs/", os.ModePerm)
	if err != nil {
		pp.Printf("Couldn't create 'logs' directory. Got: %s", err)
		panic(err)
	}
}
