package logger

import (
	json2 "encoding/json"
	"github.com/flyerxp/lib/app"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/utils/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

type errLog struct {
	ZapLog     *zap.Logger
	once       sync.Once
	errMetrics errMetrics
	isInitEd   bool
}

type errMetrics struct {
	Error []zap.Field
}

var errLogV = new(errLog)

func GetErrorLog() {
	errLogV.once.Do(func() {
		rawJSON, _ := json.Encode(config2.GetConf().App.ErrLog)
		var cfg zap.Config
		if err := json2.Unmarshal(rawJSON, &cfg); err != nil {
			log.Print(err)
		}
		errLogV.errMetrics.Error = append(errLogV.errMetrics.Error, zap.Namespace("error"))
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		errLogV.ZapLog = zap.Must(cfg.Build())
		errLogV.isInitEd = true
		_ = app.RegisterFunc("errLog", "errLog sync", func() {
			e := errLogV.ZapLog.Sync()
			if e != nil {
				log.Println(e)
			}
		})
	})
}