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

// var noticeLog zap.Logger
type AppLog struct {
	ZapLog        *zap.Logger
	once          sync.Once
	noticeMetrics noticeMetrics
	isInitEd      bool
}

// 中间件耗时
type MiddleExec struct {
	TotalExecTime int
	Count         int
	Max           int
	Avg           int
	ConnectTime   int
	ConnectCount  int
}

type MiddleExecTime struct {
	Redis    MiddleExec
	Mysql    MiddleExec
	Pulsar   MiddleExec
	Kafka    MiddleExec
	MemCache MiddleExec
	Rpc      MiddleExec
	RocketMq MiddleExec
	Elastic  MiddleExec
	Mongo    MiddleExec
}

// Log数据聚合
type noticeMetrics struct {
	Notice        []zap.Field
	Error         []zap.Field
	Warn          []zap.Field
	Middle        MiddleExecTime
	TotalExecTime int
}

func (a MiddleExec) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("total", a.TotalExecTime)
	enc.AddInt("count", a.Count)
	enc.AddInt("avg", a.Avg)
	enc.AddInt("max", a.Max)
	enc.AddInt("ConnTime", a.ConnectTime)
	enc.AddInt("ConnCount", a.ConnectCount)
	return nil
}

func (r MiddleExecTime) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	//enc.AddString("total", r.)
	if r.Redis.Count > 0 {
		_ = enc.AddObject("redis", r.Redis)
	}
	if r.MemCache.Count > 0 {
		_ = enc.AddObject("memCache", r.MemCache)
	}
	if r.Mongo.Count > 0 {
		_ = enc.AddObject("mongo", r.Mongo)
	}
	if r.Elastic.Count > 0 {
		_ = enc.AddObject("elastic", r.Elastic)
	}
	if r.Kafka.Count > 0 {
		_ = enc.AddObject("kafka", r.Kafka)
	}
	if r.Pulsar.Count > 0 {
		_ = enc.AddObject("pulsar", r.Pulsar)
	}
	if r.Rpc.Count > 0 {
		_ = enc.AddObject("rpc", r.Rpc)
	}
	if r.Mysql.Count > 0 {
		_ = enc.AddObject("mysql", r.Mysql)
	}
	if r.RocketMq.Count > 0 {
		_ = enc.AddObject("rocket", r.RocketMq)
	}
	/*zap.Inline(r.MemCache).AddTo(enc)
	zap.Inline(r.Mongo).AddTo(enc)
	zap.Inline(r.Elastic).AddTo(enc)
	zap.Inline(r.Kafka).AddTo(enc)
	zap.Inline(r.Pulsar).AddTo(enc)
	zap.Inline(r.Rpc).AddTo(enc)
	zap.Inline(r.Mysql).AddTo(enc)
	zap.Inline(r.RocketMq).AddTo(enc)*/
	return nil
}

var noticeLog = new(AppLog)

func GetNoticeLog() {
	noticeLog.once.Do(func() {
		rawJSON, _ := json.Encode(config2.GetConf().App.Logger)
		var cfg zap.Config
		if err := json2.Unmarshal(rawJSON, &cfg); err != nil {
			log.Print(err)
		}
		noticeLog.noticeMetrics.Notice = append(noticeLog.noticeMetrics.Notice, zap.Namespace("notice"))
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		noticeLog.ZapLog = zap.Must(cfg.Build())
		noticeLog.isInitEd = true
		_ = app.RegisterFunc("logger", "sync logger", func() {
			e := noticeLog.ZapLog.Sync()
			if e != nil {
				log.Println(e)
			}
		})
	})
}

func addMiddleExecTime(m *MiddleExec, t int) {
	m.Count += 1
	m.TotalExecTime += t
	m.Avg = (int)(m.TotalExecTime / m.Count)
	if t > m.Max {
		m.Max = t
	}
}

func addMiddleConnTime(m *MiddleExec, t int) {
	m.ConnectCount += 1
	m.ConnectTime += t
}
func AddMongoConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.Mongo, t)
}
func AddRedisConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.Redis, t)
}
func AddPulsarConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.Pulsar, t)
}
func AddKafkaConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.Kafka, t)
}
func AddEsConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.Elastic, t)
}
func AddRpcConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.Rpc, t)
}
func AddRocketConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.RocketMq, t)
}
func AddMysqlConnTime(t int) {
	addMiddleConnTime(&noticeLog.noticeMetrics.Middle.Mysql, t)
}
func AddMongoTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.Mongo, t)
}
func AddRedisTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.Redis, t)
}
func AddPulsarTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.Pulsar, t)
}
func AddKafkaTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.Kafka, t)
}
func AddEsTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.Elastic, t)
}
func AddRpcTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.Rpc, t)
}
func AddRocketTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.RocketMq, t)
}
func AddMysqlTime(t int) {
	addMiddleExecTime(&noticeLog.noticeMetrics.Middle.Mysql, t)
}
func SetExecTime(t int) {
	noticeLog.noticeMetrics.TotalExecTime = t
}
