package nacos

import (
	"context"
	"fmt"
	config2 "github.com/flyerxp/globalStruct/config"
	"github.com/flyerxp/lib/middleware/nacos"
	yaml "github.com/flyerxp/lib/utils/yaml"
	"testing"
)

func TestConf(t *testing.T) {
	a, e := nacos.NewClient("nacosConf", context.Background())
	r, _ := a.GetConfig(context.Background(), "zk", "zookeeper", "62c3bcf9-7948-4c26-a353-cebc0a7c9712")
	zk := new(config2.ZookeeperConf)
	_ = yaml.DecodeByBytes(r, zk)
	fmt.Println(zk, e)
}
