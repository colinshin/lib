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
	Port    string `yaml:"port" json:"port"`
	User    string `yaml:"user" json:"user"`
	Pwd     string `yaml:"pwd" json:"pwd"`
	Db      string `yaml:"db" json:"db"`
	//Ssl          string `yaml:"ssl" json:"ssl"` //true|false
	CharSet      string `yaml:"char_set" json:"char_set"`
	ReadTimeout  int    `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout" json:"write_timeout"`
	ConnTimeout  int    `yaml:"conn_timeout" json:"conn_timeout"`
	Collation    string `yaml:"collation" json:"collation"`
	MaxOpenConns int    `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns" json:"max_idle_conns"`
}
type MysqlConf struct {
	List []MidMysqlConf `yaml:"list" json:"list"`
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
