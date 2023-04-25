package logger

import "go.uber.org/zap/zapcore"

/*
	func HumanEncoderConfig() zapcore.EncoderConfig {
		cfg := testEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncodeDuration = zapcore.StringDurationEncoder
		return cfg
	}
*/
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
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
