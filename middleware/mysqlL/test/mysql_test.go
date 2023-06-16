package mysqlL

import (
	"database/sql"
	"testing"
	"time"
)

func TestConf(T *testing.T) {
	db, err := sql.Open("mysql", "test:12346@127.0.0.1/nacos")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}
