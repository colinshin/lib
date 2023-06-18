package config

import (
	"github.com/flyerxp/globalStruct/config"
	"github.com/flyerxp/lib/middleware/mysqlL"
	"github.com/flyerxp/lib/middleware/nacos"
	"github.com/flyerxp/lib/middleware/redisL"
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
	Hertz      config.Hertz        `yaml:"hertz" json:"hertz"`
	Redis      redisL.MidRedisConf `yaml:"redis" json:"redis"`
	RedisNacos nacos.NacosConf     `yaml:"redisNacos" json:"redisNacos"`
	Mysql      mysqlL.MysqlConf    `yaml:"mysql" json:"mysql"`
	Pulsar     config.PulsarConf   `yaml:"pulsar" json:"pulsar"`
	Nacos      []nacos.MidNacos    `yaml:"nacos" json:"nacos"`
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
/*func (a *Config) getRedisConf(name string) {
	if c.Redis == nil {
		err := yaml.DecodeByFile(filepath.Join(prefix, filepath.Join(env.GetEnv(), "redis.yml")), config)
		if err != nil {
			//logger.Logger.
		}
	}
}*/

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
    level: warn
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
redis:
-
  name: pubRedis
  address: [ "127.0.0.1:6379" ]
  user:
  pwd:
  master:
redisNacos:
  name: nacosConf
  did: redis
  group: redis
  ns: 62c3bcf9-7948-4c26-a353-cebc0a7c9712
nacos:
-
  name: nacosConf
  url: http://nacosconf:8848/nacos
  contextPath: /nacos
  ns: 62c3bcf9-7948-4c26-a353-cebc0a7c9712
  user: dev
  pwd: 123456
  master:
  redis:
    name: base
    address: [ "127.0.0.1:6379" ]
    user:
    pwd:
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
