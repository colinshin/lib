package redisL

import (
	"context"
	"errors"
	"github.com/flyerxp/globalStruct/config"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/nacos"
	yaml2 "github.com/flyerxp/lib/utils/yaml"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type redisClient struct {
	RedisClient cmap.ConcurrentMap[string, redis.UniversalClient]
	RedisConf   cmap.ConcurrentMap[string, config.MidRedisConf]
}

var redisEngine *redisClient

func GetEngine(name string, ctx context.Context) (redis.UniversalClient, error) {
	if redisEngine == nil {
		redisEngine = new(redisClient)
		var confList []config.MidRedisConf
		redisEngine.RedisConf = cmap.New[config.MidRedisConf]()
		redisEngine.RedisClient = cmap.New[redis.UniversalClient]()
		conf := config2.GetConf()
		confList = conf.Redis
		//本地文件中获取
		//redisEngine.Lock.Lock()
		for _, v := range confList {
			if v.Name != "" {
				redisEngine.RedisConf.Set(v.Name, v)
			}
		}
		//nacos获取
		if conf.RedisNacos.Name != "" {
			var yaml []byte
			redisList := new(config.RedisConf)
			ns, e := nacos.GetEngine(conf.RedisNacos.Name, ctx)
			if e == nil {
				yaml, e = ns.GetConfig(ctx, conf.RedisNacos.Did, conf.RedisNacos.Group, conf.RedisNacos.Ns)
				if e == nil {
					e = yaml2.DecodeByBytes(yaml, redisList)
					if e == nil {
						for _, v := range redisList.Redis {
							redisEngine.RedisConf.Set(v.Name, v)
						}
					} else {
						logger.AddError(zap.Error(errors.New("yaml conver error")))
					}
				}
			}
		}
		//redisEngine.Lock.Unlock()
	}
	e, ok := redisEngine.RedisClient.Get(name)
	if ok {
		return e, nil
	}
	o, okC := redisEngine.RedisConf.Get(name)
	if okC {
		objRedis := redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        o.Address,
			MasterName:   o.Master,
			Username:     o.User,
			Password:     o.Pwd,
			PoolTimeout:  time.Second,
			MaxIdleConns: 30,
		})
		objRedis.AddHook(HookLog{})
		redisEngine.RedisClient.Set(name, objRedis)
		return objRedis, nil
	}
	logger.AddError(zap.Error(errors.New("no find redis config " + name)))
	return nil, errors.New("no find redis config " + name)
}
func Reset() {
	redisEngine = nil
}
