# 整合Lib,避免重复造轮子

使用方法
===
* 环境变量

  GO_ENV 作为读取配置文件的目录，例如 值为test ,则读取配置文件 /conf/test/app.yml

  HOSTS:  127.0.0.1 nacosconf nacosconfredis pubmysql pubpulsar

* 配置文件
    * 默认的配置放在

      config/api.go 里，找不到配置文件，则使用这个配置，便于测试
    * 配置文件
        * [app.yml](https://github.com/flyerxp/lib/blob/main/config/test/conf/test/app.yml)

          综合的app配置
        * pulsar 的topic配置

          pulsar.yml 用来指定哪个topic，放到哪个集群 未指定clusert的，会按照 topic / 1000000 的整数，获取集群代号，按照topic_distribution 指定的配置获取集群

          topicinit.yml 是为了加速第一次producer消息，可以不配置，pulsar 的客户端建立producter 第一次会比较慢，这个配置是为了解决第一加载的问题，没有必要，不用配置

          参考 [middleware\pulsarL\test\conf\test\pulsar.yml](https://github.com/flyerxp/lib/blob/main/middleware/pulsarL/test/conf/test/pulsar.yml)
          
          参考 [topicinit](https://github.com/flyerxp/lib/blob/main/middleware/pulsarL/test/conf/test/topicinit.yml)

工具包
===
* Json工具包使用
  ```Go
     package main
     import (
	        myjson "github.com/flyerxp/lib/utils/json"
	        "strings"	     
            "fmt"
     )
     func main(){
       tmp := map[string]string{
          "a":"b",
       }
       r,e:=myjson.Encode(tmp)
       fmt.Println(string(r),e)
     }
  ```
* Yaml工具包使用







