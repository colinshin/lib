package logger

import (
	"errors"
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
	//logger.AddCtxNotice(context.Background(), zap.Int("cccc", 1111))
	logger.AddMongoTime(1)
	logger.AddMysqlTime(1)
	logger.AddPulsarTime(1)
	logger.AddKafkaTime(1)
	logger.AddRpcTime(1)
	logger.AddEsTime(1)
	logger.AddRocketTime(1)
	logger.AddRocketConnTime(1)
	logger.AddRedisConnTime(1)
	logger.AddMongoConnTime(1)
	logger.AddMysqlConnTime(1)
	logger.AddPulsarConnTime(1)
	logger.AddKafkaConnTime(1)
	logger.AddRpcConnTime(1)
	logger.AddEsConnTime(1)
	logger.AddNotice(zap.Int("cccc", 1111))
	logger.AddRedisTime(10)
	logger.SetExecTime(12)
	logger.AddNotice(zap.String("cccc", "add add add"))
	//logger.AddCtxNotice(context.Background(), zap.String("cccc", "aaaaaaaaaaaaa"))
	logger.AddError(zap.Error(errors.New("error error error")))
	//logger.AddCtxError(context.Background(), zap.Error(errors.New("error error error")))
	logger.AddError(zap.Error(errors.New("error error error")))
	//logger.AddCtxWarn(context.Background(), zap.Error(errors.New("warn warn warn")))
	logger.AddWarn(zap.Error(errors.New("warn warn warn")))
	logger.WriteErr()
	logger.WriteLine()

}
