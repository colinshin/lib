package config

type MidConf struct {
	Name    string   `yaml:"name" json:"name"`
	Address []string `yaml:"address" json:"address"`
	User    string   `yaml:"user" json:"user"`
	Pwd     string   `yaml:"pwd" json:"pwd"`
}
type ZookeeperConf struct {
	List []MidConf `yaml:"zookeeper" json:"zookeeper"`
}
type Hertz struct {
	Address         string `yaml:"address"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
}
