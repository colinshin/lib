package config

import (
	json2 "encoding/json"
	"fmt"
	"github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/utils/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestConf(t *testing.T) {
	a := config.GetConf()
	//fmt.Println(a)
	rawJSON, _ := json.Encode(a.App.Logger)
	//fmt.Print(rawJSON)
	var cfg zap.Config
	if err := json2.Unmarshal(rawJSON, &cfg); err != nil {
		fmt.Print(err)
	}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	//fmt.Printf("%#v %s", tmp, e)
	logger := zap.Must(cfg.Build())
	notice := logger.With(
		zap.Namespace("notice"),
	)
	tmpSets := make([]zap.Field, 2)
	tmpSets[0] = zap.Int("b", 2)
	tmpSets[1] = zap.Int("c", 3)
	notice.Info("notice", tmpSets...)
	//notice.Info("aaaaa", zap.Int("test2", 2))
	//a := logger.Sugar()

	//logger.Debug("xxxxxx", zap.Int("aaaaaaa", 1), zap.Any("a", a))
	defer logger.Sync()
	//defer logger.Sync()
	//var config config.Config
	//err := yaml.DecodeByFile(filepath.Join("conf", filepath.Join(env.GetEnv(), "app.yml")), &config)
	//fmt.Println(err)
	//fmt.Println(a)
	//var config config.Config
	//fmt.Println(filepath.Join("conf", filepath.Join(env.GetEnv(), "app.yml")))
	//err := yaml.DecodeByFile(filepath.Join("conf", filepath.Join(env.GetEnv(), "app.yml")), &config)
	//fmt.Println(err)
	//os.Exit(0)
}
