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

type warnLog struct {
	ZapLog      *zap.Logger
	once        sync.Once
	warnMetrics warMetrics
	isInitEd    bool
}
type warMetrics struct {
	Warn          []zap.Field
	TotalExecTime int
}

var warnLogV = new(warnLog)

func GetWarnLog() {
	warnLogV.once.Do(func() {
		rawJSON, _ := json.Encode(config2.GetConf().App.ErrLog)
		var cfg zap.Config
		if err := json2.Unmarshal(rawJSON, &cfg); err != nil {
			log.Print(err)
		}
		warnLogV.warnMetrics.Warn = append(warnLogV.warnMetrics.Warn, zap.Namespace("warn"))
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		warnLogV.ZapLog = zap.Must(cfg.Build())
		warnLogV.isInitEd = true
		_ = app.RegisterFunc("warnLog", "errLog sync", func() {
			e := warnLogV.ZapLog.Sync()
			if e != nil {
				log.Println(e)
			}
		})
	})
}
