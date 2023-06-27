# lib
go get github.com/flyerxp/lib

公共类库

Json 工具:  

import json2 "github.com/flyerxp/lib/utils/json"
data=map[string]string{"a":"b"}
json2.Encode(&data)

配置中心: nacos

redis: github.com/redis/go-redis/v9

yaml: gopkg.in/yaml.v3

log: go.uber.org/zap

协程：  ants

消息订阅分发： pulsar

mysql: github.com/jmoiron/sqlx

map cmap



HOSTS:
127.0.0.1 nacosconf nacosconfredis pubmysql pubpulsar

package main

import (
	"context"
	"fmt"
	"github.com/flyerxp/lib/app"
	"github.com/flyerxp/lib/middleware/mysqlL"
	"github.com/flyerxp/lib/middleware/pulsarL"
	"github.com/flyerxp/lib/middleware/redisL"
	"time"
)

type tmp struct {
	Id int
}

func main() {
	//time.Sleep(time.Second * 1)
	defer app.Shutdown(context.Background())
	start := time.Now()
	count := 10000
	ctx := context.Background()
	objRedis, _ := redisL.GetEngine("pubRedis", ctx)
	tmp2 := new(tmp)
	fmt.Println("github.com/flyerxp")
	fmt.Println("win11 环境，开始了 ")
	for i := 0; i <= count; i++ {
		objRedis.Get(context.Background(), "a")
	}
	fmt.Printf("redis 读取 10000次耗时 %d 毫秒\n", time.Since(start).Milliseconds())
	start = time.Now()
	mysql, _ := mysqlL.GetEngine("pubMysql", context.Background())
	for i := 0; i <= count; i++ {
		err := mysql.GetDb().Get(tmp2, `select id from config_info limit 1`)
		if err != nil {
			fmt.Println(tmp2, err)
		}
	}
	fmt.Printf("mysql 数据库读取 10000次耗时 %d 毫秒\n", time.Since(start).Milliseconds())
	start = time.Now()
	for i := 0; i <= count; i++ {
		pulsarL.Producer(&pulsarL.OutMessage{
			Topic:      0,
			TopicStr:   "test",
			Content:    "太牛了",
			Properties: map[string]string{"a": "b"},
			Delay:      0,
		}, ctx)
	}
	fmt.Printf("pulsar 发消息 10000次耗时 %d 毫秒\n", time.Since(start).Milliseconds())
}





