package redisL

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
