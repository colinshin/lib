package logger

import (
	"go.uber.org/zap"
)

var libLog zap.Logger

func setLogger() {

}
func getLogger() {

}

/*
	func HumanEncoderConfig() zapcore.EncoderConfig {
		cfg := testEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncodeDuration = zapcore.StringDurationEncoder
		return cfg
	}
*/
