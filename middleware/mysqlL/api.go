package mysqlL

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"sync"
)

type MidConf struct {
	Name    string   `yaml:"name" json:"name"`
	Address []string `yaml:"address" json:"address"`
	User    string   `yaml:"user" json:"user"`
	Pwd     string   `yaml:"pwd" json:"pwd"`
}
type MysqlConf struct {
	List []MidConf `yaml:"list" json:"list"`
}
type MysqlClient struct {
	MysqlClient cmap.ConcurrentMap[string, *sync.Pool]
	MysqlConf   cmap.ConcurrentMap[string, MysqlConf]
}

func GetEngine() {

}
