package mysqlL

import (
	"context"
	"fmt"
	"github.com/flyerxp/lib/middleware/mysqlL"
	"testing"
	"time"
)

type tmp struct {
	Id int
}

func TestConf(T *testing.T) {
	tmp := new(tmp)
	count := 10000
	start := time.Now()
	mysql, _ := mysqlL.GetEngine("pubMysql", context.Background())
	for i := 0; i <= count; i++ {
		err := mysql.GetDb().Get(tmp, `select id from config_info limit 1`)
		if err != nil {
			fmt.Println(tmp, err)
		}
	}
	fmt.Printf("mysql 数据库读取 10000次耗时 %d 毫秒\n", time.Since(start).Milliseconds())
}
