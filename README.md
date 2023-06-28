# 整合Lib,避免重复造轮子

使用方法
===
* 依赖
  
  nacos 配置中心，配置中心的数据，会在redis生成缓存，更新后，需要清理缓存，参考示例 nacos [监听事件订阅](https://github.com/nacos-group/nacos-sdk-go/blob/master/README_CN.md),订阅到事件后，删除缓存
  
  Lib 对外提供删除key的方法
 
  ```Go
  package main
  /*
  client, e := nacos.GetEngine("nacosConf", context.Background())
				if e != nil {
					logger.AddError(zap.Error(e))
				}
				key := client.DeleteCache(context.Background(), dataId, group, syncConf.Ns)
  */
  ```

  redis 缓存
  

* 环境变量

  GO_ENV 作为读取配置文件的目录，例如 值为test ,则读取配置文件 /conf/test/app.yml

  HOSTS:  127.0.0.1 nacosconf nacosconfredis pubmysql pubpulsar pubredis

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
  /*
  import (
        "strings"	     
        "fmt"
        myjson "github.com/flyerxp/lib/utils/json"
     )
     func main(){
       tmp := map[string]string{
          "a":"b",
       }
       r,e:=myjson.Encode(tmp)
       fmt.Println(string(r),e)
  } */
  ```
* Yaml工具包使用

  ```Go
  package main
  /*
  import (
    "fmt"
    myyml "github.com/flyerxp/lib/utils/yaml"
  )  

  func main() {
      var defaultConfig = []byte(`
      a: b
      `)
  tmp := map[string]string{}
  //myyml.DecodeByFile("app.yml", &tmp)
  myyml.DecodeByBytes(defaultConfig, tmp)
  fmt.Println(string(defaultConfig))
  }*/
  ```

* 中间件使用

    ```Go
    package main
    /*
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
      }*/

    ```
    [测试结果](https://github.com/flyerxp/lib/blob/main/doc/image/test.png)

  






