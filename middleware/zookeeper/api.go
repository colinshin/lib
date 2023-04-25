package zookeeper

import (
	config "github.com/flyerxp/globalStruct/config"
	"github.com/flyerxp/lib/utils/env"
	"os"
	"path/filepath"
	"sync"
)

func init() {

	/*&sync.Pool{
		localSize: 10,
		New: func() interface{} {
			fmt.Println("creating a new person")
			c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*2)
			if err != nil {
				panic(err)
			}
			return c
		},
	}*/
}

type Engine struct {
	Conf   config.ZookeeperConf
	ZkPool map[string]*sync.Pool
	Once   sync.Once
}

func (z *Engine) SetLogger(f func()) {

}
func (z *Engine) New(cluster string) {
	z.Once.Do(func() {
		prefix := "conf"
		confFileRelPath := filepath.Join(prefix, filepath.Join(env.GetEnv(), "zookeeper.yml"))
		_, err := os.Stat(confFileRelPath)
		if err != nil && os.IsNotExist(err) {
			//log.Logger.Writer()
		} else {

		}
	})
}
