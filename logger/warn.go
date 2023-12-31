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

func getWarnLog() {
	warnLogV.once.Do(func() {
		rawJSON, _ := json.Encode(config2.GetConf().App.ErrLog)
		var cfg zap.Config
		if err := json2.Unmarshal(rawJSON, &cfg); err != nil {
			log.Print(err)
		}
		cfg.OutputPaths = GetPath(cfg.OutputPaths, "warn")
		warnLogV.warnMetrics.Warn = append(warnLogV.warnMetrics.Warn, zap.Namespace("warn"))
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		warnLogV.ZapLog = zap.Must(cfg.Build())
		warnLogV.isInitEd = true
		RegistermakeFileEvent(Event{"error", func() {
			warnLogV = new(warnLog)
			getWarnLog()
		}})
		RegisterReset(Event{"notice", func() {
			warnLogV.warnMetrics.Warn = make([]zap.Field, 1, 10)
			warnLogV.warnMetrics.Warn[0] = zap.Namespace("warn")
		}})
		if len(cfg.OutputPaths) > 0 {
			_ = app.RegisterFunc("warnLog", "errLog sync", func() {
				_ = warnLogV.ZapLog.Sync()
			})
		}
	})
}
