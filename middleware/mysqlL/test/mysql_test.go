package mysqlL

import (
	"context"
	"fmt"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/mysqlL"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

type tmp struct {
	id int
}

func TestConf(T *testing.T) {

	defer logger.WriteLine()
	r, err := mysqlL.GetEngine("pubMysql", context.Background())

	conn := r.Get().(*sqlx.DB)

	var tmp2 tmp
	//conn.Get(&tmp2, "select *  from config_info limit 1")
	s := time.Now()
	err = conn.Get(&tmp2, `select * from config_info limit 1`)
	fmt.Println(time.Since(s).Milliseconds(), "我是总耗时")
	fmt.Println(tmp2, err)
	//r, e := conn.Query("select * from config_info")
	//fmt.Println(e, r)
	//conn.Ping()
	//fmt.Println(conn)

	time.Sleep(time.Second)
	//logger.WriteLine()
}
