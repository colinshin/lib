package config

import (
	"fmt"
	"github.com/flyerxp/lib/config"
	"testing"
)

func TestConf(t *testing.T) {
	a := config.GetConf().RedisNacos
	fmt.Println(a)
}
