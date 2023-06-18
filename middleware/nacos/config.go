package nacos

import "github.com/flyerxp/lib/middleware/redisL"

type NacosConf struct {
	Name  string `yaml:"name" json:"name"`
	Did   string `yaml:"did" json:"did"`
	Group string `yaml:"group" json:"group"`
	Ns    string `yaml:"ns" json:"ns"`
}
type MidNacos struct {
	Name        string              `yaml:"name" json:"name"`
	Url         string              `yaml:"url" json:"url"`
	ContextPath string              `yaml:"contextPath" json:"contextPath"`
	Ns          string              `yaml:"ns" json:"ns"`
	Redis       redisL.MidRedisConf `yaml:"redis" json:"redis"`
	User        string              `yaml:"user" json:"user"`
	Pwd         string              `yaml:"pwd" json:"pwd"`
}
type Nacos struct {
	List []MidNacos `yaml:"nacos"`
}
