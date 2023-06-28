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




