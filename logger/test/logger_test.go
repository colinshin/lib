package logger

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flyerxp/lib/logger"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	/*dynamicLevel := zap.NewAtomicLevel()
	dynamicLevel.SetLevel(zap.DebugLevel)*/
	l := hertzzap.NewLogger(
		hertzzap.WithCores([]hertzzap.CoreConfig{
			{
				Enc: zapcore.NewConsoleEncoder(logger.EncoderConfig()),
				Ws:  zapcore.AddSync(os.Stdout),
				Lvl: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			},
			{
				Enc: zapcore.NewJSONEncoder(logger.EncoderConfig()),
				Ws:  getWriteSyncer("hertz.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.DebugLevel),
				/*Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.DebugLevel
				}),*/
			},
		}...),
	)
	defer l.Sync()
	hlog.SetLogger(l)
	hlog.Notice("notice log")
	hlog.Notice("notice log2")
	hlog.Infof("hello %s", "hertz")
	hlog.Info("hertz")
	hlog.Warn("hertz")
	hlog.Error("error")
	hlog.Debugf("xxxxxxxxxxxxxxx")
}
func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    1000,
		MaxBackups: 5,
		MaxAge:     48,
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
func Test2Encode(t *testing.T) {
	logger.GetNoticeLog()
	logger.AddNotice(zap.Int("cccc", 1111))
	logger.AddRedisTime(10)
	logger.AddNotice(zap.String("cccc", "aaaaaaaaaaaaa"))
	logger.GetErrorLog().Info("error", zap.Int("aaaaa", 1))
}
