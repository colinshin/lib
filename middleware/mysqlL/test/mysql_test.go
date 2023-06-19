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

func TestConf(T *testing.T) {
	s := time.Now()
	defer logger.WriteLine()
	r, _ := mysqlL.GetEngine("pubMysql", context.Background())
	//fmt.Println(err)
	_ = r.Get().(*sqlx.DB)
	//fmt.Println(conn)
	fmt.Println(time.Since(s).Milliseconds())
	time.Sleep(time.Second)
	//logger.WriteLine()
}
