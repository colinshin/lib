package redisL

import (
	"context"
	"github.com/flyerxp/lib/middleware/redisL"
	"testing"
)

func TestConf(t *testing.T) {
	r, _ := redisL.GetEngine("pubRedis", context.Background())
	r.Get(context.Background(), "aaaa")
}
