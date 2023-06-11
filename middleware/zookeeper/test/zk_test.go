package test

import (
	"fmt"
	"github.com/flyerxp/lib/middleware/zookeeper"
	"testing"
)

func TestConf(t *testing.T) {
	zkConn := zookeeper.New("centerZk")
	//defer zookeeper.PutConn("centerZk", zkConn)
	a, s, e := zkConn.Get("/")
	fmt.Printf("%#v", a)
	fmt.Printf("%#v", s)
	fmt.Printf("%#v", e)
	//zkConn.Close()
	zookeeper.PutConn("centerZk", zkConn)
	//time.Sleep(time.Second * 20)
	zkConn = zookeeper.New("centerZk")
	a, s, e = zkConn.Get("/")
	fmt.Println("===================================")
	fmt.Printf("%#v", a)
	fmt.Printf("%#v", s)
	fmt.Printf("%#v", e)
	zookeeper.Reset()
}
