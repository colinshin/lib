package config

type MidConf struct {
	Name    string   `yaml:"name" json:"name"`
	Address []string `yaml:"address" json:"address"`
	User    string   `yaml:"user" json:"user"`
	Pwd     string   `yaml:"pwd" json:"pwd"`
}
type MidMysqlConf struct {
	Name    string `yaml:"name" json:"name"`
	Address string `yaml:"address" json:"address"`
	User    string `yaml:"user" json:"user"`
	Pwd     string `yaml:"pwd" json:"pwd"`
	Db      string `yaml:"db" json:"db"`
	Ssl     string `yaml:"ssl" json:"ssl"` //disable|enable
}
type MysqlConf struct {
	List []MidConf `yaml:"list" json:"list"`
}
type ZookeeperConf struct {
	List []MidConf `yaml:"zookeeper" json:"zookeeper"`
}
type NacosConf struct {
	Name  string `yaml:"name" json:"name"`
	Did   string `yaml:"did" json:"did"`
	Group string `yaml:"group" json:"group"`
	Ns    string `yaml:"ns" json:"ns"`
}
type MidNacos struct {
	Name        string       `yaml:"name" json:"name"`
	Url         string       `yaml:"url" json:"url"`
	ContextPath string       `yaml:"contextPath" json:"contextPath"`
	Ns          string       `yaml:"ns" json:"ns"`
	Redis       MidRedisConf `yaml:"redis" json:"redis"`
	User        string       `yaml:"user" json:"user"`
	Pwd         string       `yaml:"pwd" json:"pwd"`
}
type Nacos struct {
	List []MidNacos `yaml:"nacos"`
}

type MidRedisConf struct {
	Name    string   `yaml:"name" json:"name"`
	Address []string `yaml:"address" json:"address"`
	User    string   `yaml:"user" json:"user"`
	Pwd     string   `yaml:"pwd" json:"pwd"`
	Master  string   `yaml:"master" json:"master"` //哨兵模式使用，写masterName
}
type RedisConf struct {
	Redis []MidRedisConf `yaml:"redis" json:"redis"`
}

type ElasticConf struct {
	List []MidConf `yaml:"elastic" json:"elastic"`
}
