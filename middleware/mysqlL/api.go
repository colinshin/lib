package mysqlL

import (
	"context"
	"errors"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/jmoiron/sqlx"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
)

type MysqlClient struct {
	MysqlClient cmap.ConcurrentMap[string, *sync.Pool]
	MysqlConf   cmap.ConcurrentMap[string, config2.MysqlConf]
}

var redisClient redis.UniversalClient

func GetEngine(name string, ctx context.Context) (*sync.Pool, error) {
	for _, v := range config2.GetConf().Mysql {
		if v.Name == name {
			return newClient(v, ctx), nil
		}
	}
	logger.AddError(zap.Error(errors.New("nacos conf no find " + name)))
	return nil, errors.New("nacos conf no find " + name)
}
func newClient(o config2.MidMysqlConf, ctx context.Context) *sync.Pool {
	c := &sync.Pool{
		New: func() any {
			var o config2.MidMysqlConf
			n := sqlx.MustConnect("postgres", "user="+o.User+" dbname=(tpc:"+o.Address+")"+o.Name+" sslmode="+o.Ssl)
			return n
		},
	}
	return c
}
