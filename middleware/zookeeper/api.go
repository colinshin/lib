package zookeeper

import (
	"errors"
	config "github.com/flyerxp/globalStruct/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/utils/env"
	yaml "github.com/flyerxp/lib/utils/yaml"
	"github.com/go-zookeeper/zk"
	"go.uber.org/zap"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

var zkEngine = new(Engine)

type Engine struct {
	Conf       *config.ZookeeperConf
	ZkPool     map[string]*sync.Pool
	Once       sync.Once
	ReloadNums int32
}

// 获取资源
func New(cluster string) *zk.Conn {
	zkEngine.Once.Do(func() {
		initConfig()
	})
	zkPool, ok := zkEngine.ZkPool[cluster]
	if !ok {
		for i, v := range zkEngine.Conf.List {
			if zkEngine.Conf.List[i].Name == cluster {
				zkEngine.ZkPool[cluster] = &sync.Pool{
					New: func() any {
						//fmt.Println("\n", "i im createing----------------------------------")
						c, _, e := zk.Connect(v.Address, time.Second)
						if e != nil {
							logger.GetErrorLog().Error("zookeeper "+cluster, zap.Error(e))
						}
						return c
					},
				}
				zkPool, ok = zkEngine.ZkPool[cluster]

			}
		}
	}
	if ok {
		return zkPool.Get().(*zk.Conn)
	} else {
		if zkEngine.ReloadNums <= 1 {
			initConfig()
			atomic.AddInt32(&zkEngine.ReloadNums, 1)
		} else if zkEngine.ReloadNums > 1 {
			time.AfterFunc(time.Second*3, func() {
				initConfig()
				if zkEngine.ReloadNums > 6 {
					zkEngine.ReloadNums = 0
				}
			})
		}
		logger.GetErrorLog().Error("zookeeper "+cluster, zap.Error(errors.New("no find config")))
	}
	return nil
}

// 归还资源,如果不归还资源，则资源无法被重复利用
func PutConn(c string, conn *zk.Conn) {
	zkPool, ok := zkEngine.ZkPool[c]
	if ok {
		zkPool.Put(conn)
	}
}
func initConfig() {
	//fmt.Println("i reload")
	prefix := "conf"
	conf := new(config.ZookeeperConf)
	err := yaml.DecodeByFile(filepath.Join(prefix, filepath.Join(env.GetEnv(), "zookeeper.yml")), conf)
	if err != nil {
		logger.GetErrorLog().Warn("zookeeper error", zap.Error(err))
	} else {
		zkEngine.Conf = conf
	}
	zkEngine.ZkPool = make(map[string]*sync.Pool)
}
func Reset() {
	initConfig()
}
