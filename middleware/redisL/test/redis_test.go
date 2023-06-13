package redisL

import (
	"context"
	"github.com/flyerxp/lib/app"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/redisL"
	"go.uber.org/zap"
	"testing"
)

func TestConf(t *testing.T) {

	logger.AddNotice(zap.String("a", "cccccccccccccccc"))

	r, _ := redisL.GetEngine("pubRedis", context.Background())
	big1 := logger.StartTime("beg1")
	l := logger.StartTime("redis-read")
	r.Get(context.Background(), "aaaa")
	l.Stop()
	l2 := logger.StartTime("redis-read")
	//time.Sleep(time.Second)
	l2.Stop()
	big1.Stop()
	logger.WriteLine()
	app.Shutdown(context.Background())
	redisL.Reset()
}
