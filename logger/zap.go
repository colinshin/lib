package logger

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zap log 日志通用格式
func EncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05.000Z0700"),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
func GetWriteSyncer(file string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    50,
		MaxBackups: 10,
		MaxAge:     48,
		Compress:   true,
		LocalTime:  true,
	}
}
