package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ConsoleLogger *zap.Logger
	FileLogger    *zap.Logger
	LogFile       *lumberjack.Logger
)

func SetupLoggers(logPath string) error {
	LogFile = &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    50, // Megabytes
		MaxBackups: 30,  // Max number of files
		Compress:   true,
	}

	// Only dev environment
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(
		"2006-01-02 15:04:05")
	consoleEncoderConfig.TimeKey = "time"
	consoleEncoderConfig.LevelKey = "level"
	consoleEncoderConfig.CallerKey = "caller"
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// Dev and prod environments
	jsonEncoderConfig := zap.NewProductionEncoderConfig()
	jsonEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoderConfig.TimeKey = "timestamp"
	jsonEncoderConfig.CallerKey = "caller"
	jsonEncoder := zapcore.NewJSONEncoder(jsonEncoderConfig)

	level, err := levelFromEnv()
	if err != nil {
		return err
	}

	consoleCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level),
	)

	fileCore := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(LogFile), level),
	)

	ConsoleLogger = zap.New(consoleCore, zap.AddCaller())
	FileLogger = zap.New(fileCore, zap.AddCaller())
	return nil
}

// levelFromEnv reads LOG_LEVEL (debug|info|warn|error), defaulting to info.
func levelFromEnv() (zapcore.Level, error) {
	switch os.Getenv("LOG_LEVEL") {
	case "", "info":
		return zapcore.InfoLevel, nil
	case "debug":
		return zapcore.DebugLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf(
			"invalid LOG_LEVEL %q: must be debug, info, warn or error",
			os.Getenv("LOG_LEVEL"))
	}
}
