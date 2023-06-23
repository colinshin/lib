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

type pulsarLog struct {
	ZapLog      *zap.Logger
	once        sync.Once
	warnMetrics pulsarMetrics
	isInitEd    bool
}
type pulsarMetrics struct {
	Warn          []zap.Field
	TotalExecTime int
}

var pulsarLogV = new(pulsarLog)

func getpulsarLog() {
	pulsarLogV.once.Do(func() {
		rawJSON, _ := json.Encode(config2.GetConf().App.ErrLog)
		var cfg zap.Config
		if err := json2.Unmarshal(rawJSON, &cfg); err != nil {
			log.Print(err)
		}
		pulsarLogV.warnMetrics.Warn = append(pulsarLogV.warnMetrics.Warn, zap.Namespace("warn"))
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		pulsarLogV.ZapLog = zap.Must(cfg.Build())
		pulsarLogV.isInitEd = true
		_ = app.RegisterFunc("pulsarLog", "errLog sync", func() {
			e := pulsarLogV.ZapLog.Sync()
			if e != nil {
				log.Println(e)
			}
		})
	})
}
