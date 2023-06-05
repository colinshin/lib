package config

import (
	"github.com/flyerxp/globalStruct/config"
	"github.com/flyerxp/lib/utils/env"
	"github.com/flyerxp/lib/utils/json"
	yaml "github.com/flyerxp/lib/utils/yaml"
	_ "go.uber.org/zap/zapcore"
	"log"
	"path/filepath"
	"sync"
)

var (
	prefix = "conf"
	conf   *Config
	once   sync.Once
)

type zapConfig struct {
	Level            string            `yaml:"level" json:"level"`
	Encoding         string            `yaml:"encoding" json:"encoding"`
	OutputPaths      []string          `yaml:"outputPaths" json:"outputPaths"`
	ErrorOutputPaths []string          `yaml:"errorOutputPaths" json:"errorOutputPaths"`
	InitialFields    map[string]string `yaml:"initialFields" json:"initialFields"`
	EncoderConfig    map[string]string `yaml:"encoderConfig" json:"encoderConfig"`
}
type Config struct {
	Env string `yaml:"env" json:"env"`
	App struct {
		Name        string    `yaml:"name" json:"name"`
		Type        string    `yaml:"type" json:"type"`
		Logger      zapConfig `yaml:"logger" json:"logger"`
		ErrLog      zapConfig `yaml:"errlog" json:"errlog"`
		ConfStorage bool      `yaml:"confStorage" json:"confStorage"`
	}
	Hertz  config.Hertz        `yaml:"hertz" json:"hertz"`
	Redis  []config.RedisConf  `yaml:"redis" json:"redis"`
	Mysql  []config.MysqlConf  `yaml:"mysql" json:"mysql"`
	Pulsar []config.PulsarConf `yaml:"pulsar" json:"pulsar"`
}

func (c *Config) String() string {
	b, e := json.Encode(c)
	if e != nil {
		log.Fatalf("config josn error %s", e)
	}
	return string(b)
}

func GetConf() *Config {
	once.Do(initConf)
	return conf
}

// func (a *Config) getLoggerConf() zap.Config {
// return a.App.Logger
// }
func (a *Config) getRedisConf(name string) {
	/*if c.Redis == nil {
		err := yaml.DecodeByFile(filepath.Join(prefix, filepath.Join(env.GetEnv(), "redis.yml")), config)
		if err != nil {
			//logger.Logger.
		}
	}*/
}

var defaultConfig = []byte(`
env: test
app:
  name: Webhook
  type: web
  logger:
    level: debug
    encoding: json
    outputPaths:
      - stdout
      #- logs/webhook
    errorOutputPaths:
      - stderr
    initialFields:
      app: Webhook
    encoderConfig:
      #messageKey: msg
      levelKey: level
      nameKey: name
      TimeKey: time
      #CallerKey: caller
      #FunctionKey: func
      StacktraceKey: stacktrace
      LineEnding: "\n"
  errlog:
    level: debug
    encoding: json
    outputPaths:
      - stdout
      #- logs/webhook
    errorOutputPaths:
      - stderr
    initialFields:
      app: Webhook
    encoderConfig:
      #messageKey: msg
      levelKey: level
      nameKey: name
      TimeKey: time
      CallerKey: caller
      FunctionKey: func
      StacktraceKey: stacktrace
      LineEnding: "\n"
`)

func initConf() {
	conf = new(Config)
	err := yaml.DecodeByFile(filepath.Join(prefix, filepath.Join(env.GetEnv(), "app.yml")), conf)
	if err != nil {
		log.Printf("default conf no find %v", err)
		log.Print("use default config")
		err = yaml.DecodeByBytes(defaultConfig, conf)
		if err != nil {
			log.Printf("default config error", err)
		}
	}
}
