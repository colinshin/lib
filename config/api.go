package config

import (
	"github.com/flyerxp/globalStruct/config"
	"github.com/flyerxp/lib/utils/env"
	yaml "github.com/flyerxp/lib/utils/yaml"
	"log"
	"path/filepath"
	"sync"
)

var (
	prefix = "conf"
	conf   *Config
	once   sync.Once
)

type Config struct {
	Env string
	App struct {
		Name        string `yaml:"name" json:"name"`
		Type        string `yaml:"type" json:"type"`
		Logger      string `yaml:"logger" json:"logger"`
		ConfStorage bool   `yaml:"confStorage" json:"confStorage"`
	}
	Hertz  config.Hertz      `yaml:"hertz" json:"hertz"`
	Redis  config.RedisConf  `yaml:"redis" json:"redis"`
	Mysql  config.MysqlConf  `yaml:"mysql" json:"mysql"`
	Pulsar config.PulsarConf `yaml:"pulsar" json:"pulsar"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}
func (c *Config) getRedisConf(name string) {
	/*if c.Redis == nil {
		err := yaml.DecodeByFile(filepath.Join(prefix, filepath.Join(env.GetEnv(), "redis.yml")), config)
		if err != nil {
			//logger.Logger.
		}
	}*/
}

func initConf() {
	config := new(Config)
	err := yaml.DecodeByFile(filepath.Join(prefix, filepath.Join(env.GetEnv(), "conf.yml")), config)
	if err != nil {
		log.Printf("default conf no find %v", err)
	}
	conf.Env = env.GetEnv()
}
