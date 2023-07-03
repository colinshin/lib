package redisL

import (
	"context"
	"errors"
	"github.com/flyerxp/lib/app"
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
	RedisClient cmap.ConcurrentMap[string, *RedisC]
	RedisConf   cmap.ConcurrentMap[string, config2.MidRedisConf]
}

var RedisEngine *redisClient

type RedisC struct {
	C redis.UniversalClient
}

func GetEngine(name string, ctx context.Context) (*RedisC, error) {
	if RedisEngine == nil {
		RedisEngine = new(redisClient)
		var confList []config2.MidRedisConf
		RedisEngine.RedisConf = cmap.New[config2.MidRedisConf]()
		RedisEngine.RedisClient = cmap.New[*RedisC]()
		conf := config2.GetConf()
		confList = conf.Redis
		//本地文件中获取
		//RedisEngine.Lock.Lock()
		for _, v := range confList {
			if v.Name != "" {
				RedisEngine.RedisConf.Set(v.Name, v)
			}
		}
		//nacos获取
		if conf.RedisNacos.Name != "" {
			var yaml []byte
			redisList := new(config2.RedisConf)
			ns, e := nacos.GetEngine(conf.RedisNacos.Name, ctx)
			if e == nil {
				yaml, e = ns.GetConfig(ctx, conf.RedisNacos.Did, conf.RedisNacos.Group, conf.RedisNacos.Ns)
				if e == nil {
					e = yaml2.DecodeByBytes(yaml, redisList)
					if e == nil {
						for _, v := range redisList.Redis {
							RedisEngine.RedisConf.Set(v.Name, v)
						}
					} else {
						logger.AddError(zap.Error(errors.New("yaml conver error")))
					}
				}
			}
		}
		_ = app.RegisterFunc("redis", "redis close", func() {
			RedisEngine.Reset()
		})
		//RedisEngine.Lock.Unlock()
	}
	e, ok := RedisEngine.RedisClient.Get(name)
	if ok {
		return e, nil
	}
	o, okC := RedisEngine.RedisConf.Get(name)
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
		objRedisC := new(RedisC)
		objRedisC.C = objRedis
		RedisEngine.RedisClient.Set(name, objRedisC)
		return objRedisC, nil
	}
	logger.AddError(zap.Error(errors.New("no find redis config " + name)))
	return nil, errors.New("no find redis config " + name)
}
func (r *redisClient) Reset() {
	for _, v := range RedisEngine.RedisClient.Items() {
		_ = v.C.Close()
	}
	RedisEngine = nil
}
