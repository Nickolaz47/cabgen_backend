package logging

import (
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

func SetupLoggers(logPath string) {
	LogFile = &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    50, // Megabytes
		MaxBackups: 30,  // Max number of files
		Compress:   true,
	}

	// Only dev environment
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
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

	logLevel := zapcore.DebugLevel

	consoleCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel),
	)

	fileCore := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(LogFile), logLevel),
	)

	ConsoleLogger = zap.New(consoleCore, zap.AddCaller())
	FileLogger = zap.New(fileCore, zap.AddCaller())
}
