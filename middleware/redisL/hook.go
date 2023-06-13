package redisL

import (
	"context"
	"github.com/flyerxp/lib/logger"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

type HookLog struct{}

func (HookLog) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		t := time.Now()
		l := logger.StartTime(addr)
		c, e := next(ctx, network, addr)
		l.Stop()
		logger.AddRedisConnTime(int(time.Since(t).Milliseconds()))
		return c, e
	}
}
func (HookLog) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		t := time.Now()
		e := next(ctx, cmd)
		logger.AddRedisTime(int(time.Since(t).Milliseconds()))
		return e
	}
}
func (HookLog) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		t := time.Now()
		c := next(ctx, cmds)
		logger.AddRedisTime(int(time.Since(t).Milliseconds()))
		return c
	}
}
