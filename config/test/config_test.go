package config

import (
	"fmt"
	"github.com/flyerxp/lib/config"
	"testing"
)

func TestConf(t *testing.T) {
	a := config.GetConf().Nacos
	fmt.Println(a)
}
