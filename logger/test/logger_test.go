package logger

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flyerxp/lib/logger"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	dynamicLevel := zap.NewAtomicLevel()
	dynamicLevel.SetLevel(zap.DebugLevel)
	l := hertzzap.NewLogger(
		hertzzap.WithCores([]hertzzap.CoreConfig{
			{
				Enc: zapcore.NewConsoleEncoder(logger.EncoderConfig()),
				Ws:  zapcore.AddSync(os.Stdout),
				Lvl: dynamicLevel,
			},
			{
				Enc: zapcore.NewJSONEncoder(logger.EncoderConfig()),
				Ws:  getWriteSyncer("all_log.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			},
			{
				Enc: zapcore.NewJSONEncoder(logger.EncoderConfig()),
				Ws:  getWriteSyncer("debug_log.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.DebugLevel
				}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(logger.EncoderConfig()),
				Ws:  getWriteSyncer("info_log.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.InfoLevel
				}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(logger.EncoderConfig()),
				Ws:  getWriteSyncer("warn_log.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.WarnLevel
				}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(logger.EncoderConfig()),
				Ws:  getWriteSyncer("error_log.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.ErrorLevel
				}))),
			},
		}...),
	)
	defer l.Sync()
	hlog.SetLogger(l)
	hlog.Infof("hello %s", "hertz")
	hlog.Info("hertz")
	hlog.Warn("hertz")
	hlog.Debugf("xxxxxxxxxxxxxxx")
}
func getWriteSyncer(file string) zapcore.WriteSyncer {
	return zapcore.AddSync(getWriteSyncer(file))
}
