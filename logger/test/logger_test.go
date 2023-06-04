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

func aTestEncode(t *testing.T) {
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

var libLog2 = new(logger.AppLog)

func Test2Encode(t *testing.T) {
	logger.GetNoticeLog()

	logger.AddNotice(zap.Int("cccc", 1111))
	logger.AddRedisTime(10)
	logger.AddNotice(zap.String("cccc", "aaaaaaaaaaaaa"))
	//cfg.InitialFields["notice"] = libLog2.LogMetrics.Notice
	//cfg.InitialFields["execTime"] = libLog2.LogMetrics.TotalExecTime
	//cfg.InitialFields["middle"] = libLog2.LogMetrics.Middle
	/*cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	_, e := cfg.Build()
	fmt.Printf("%#v", e)*/
	logger.WriteLine()
	//logger.WriteLine()
	//logger.WriteLine()
	//libLog2.ZapLog = zap.Must(cfg.Build())
	//logger.AddRedisTime(10)
	/*libLog2.ZapLog.With(zap.Namespace("notice"),
	zap.Int("counter", 1)).With(zap.Namespace("notice"), zap.Int("bbbbbb", 1))*/
	//libLog2.ZapLog.Info("execTime", zap.Int("execTime", 1))

}
