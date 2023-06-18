package elastic

type MidConf struct {
	Name    string   `yaml:"name" json:"name"`
	Address []string `yaml:"address" json:"address"`
	User    string   `yaml:"user" json:"user"`
	Pwd     string   `yaml:"pwd" json:"pwd"`
}

type ElasticConf struct {
	List []MidConf `yaml:"elastic" json:"elastic"`
}
