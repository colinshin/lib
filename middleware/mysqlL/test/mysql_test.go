package mysqlL

import (
	"context"
	"fmt"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/mysqlL"
	"testing"
	"time"
)

type tmp struct {
	id int
}

func TestConf(T *testing.T) {

	r, err := mysqlL.GetEngine("pubMysql", context.Background())
	conn := r.GetDb()
	time.Sleep(time.Second * 1)
	//conn.Ping()
	/*go func() {
		_ = conn.Ping()
	}()
	*/
	var tmp2 tmp
	//conn.Ping()
	//conn.Get(&tmp2, "select *  from config_info limit 1")
	s := time.Now()
	/*for {
		err = conn.Get(&tmp2, `select * from config_info limit 1`)
		time.Sleep(time.Second * 5)
		fmt.Println("oooooooooooooooooo", err)
	}*/
	err = conn.Get(&tmp2, `select * from config_info limit 1`)

	//_, err = conn.Exec(`set names utf8`)
	fmt.Println(time.Since(s).Milliseconds(), "我是总耗时")
	fmt.Println(tmp2, err)
	logger.WriteLine()
	//r, e := conn.Query("select * from config_info")
	//fmt.Println(e, r)
	//conn.Ping()
	//fmt.Println(conn)

	time.Sleep(time.Second)
	//logger.WriteLine()
}
