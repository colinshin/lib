package mysqlL

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/nacos"
	yaml2 "github.com/flyerxp/lib/utils/yaml"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/qustavo/sqlhooks/v2"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type MysqlClient struct {
	MysqlClient cmap.ConcurrentMap[string, *sync.Pool]
	MysqlConf   cmap.ConcurrentMap[string, config2.MidMysqlConf]
}

var mysqlEngine *MysqlClient

func GetEngine(name string, ctx context.Context) (*sync.Pool, error) {

	if mysqlEngine == nil {
		mysqlEngine = new(MysqlClient)
		var confList []config2.MidMysqlConf
		mysqlEngine.MysqlConf = cmap.New[config2.MidMysqlConf]()
		mysqlEngine.MysqlClient = cmap.New[*sync.Pool]()
		conf := config2.GetConf()
		confList = conf.Mysql
		//本地文件中获取
		for _, v := range confList {
			if v.Name != "" {
				mysqlEngine.MysqlConf.Set(v.Name, v)
			}
		}

		//nacos获取
		if conf.MysqlNacos.Name != "" {
			var yaml []byte
			mysqlList := new(config2.MysqlConf)
			ns, e := nacos.GetEngine(conf.MysqlNacos.Name, ctx)
			if e == nil {
				yaml, e = ns.GetConfig(ctx, conf.MysqlNacos.Did, conf.MysqlNacos.Group, conf.MysqlNacos.Ns)
				if e == nil {
					e = yaml2.DecodeByBytes(yaml, mysqlList)
					if e == nil {
						for _, v := range mysqlList.List {
							mysqlEngine.MysqlConf.Set(v.Name, v)
						}
					} else {
						logger.AddError(zap.Error(errors.New("yaml conver error")))
					}
				}
			}
		}
	}

	e, ok := mysqlEngine.MysqlClient.Get(name)
	if ok {
		return e, nil
	}
	o, okC := mysqlEngine.MysqlConf.Get(name)
	if okC {
		objMysql := newClient(o)
		mysqlEngine.MysqlClient.Set(name, objMysql)

		return objMysql, nil
	}
	logger.AddError(zap.Error(errors.New("no find mysql config " + name)))
	return nil, errors.New("no find mysql config " + name)
}

// https://github.com/golang-migrate/migrate/blob/master/database/mysql/README.md
func newClient(o config2.MidMysqlConf) *sync.Pool {
	c := &sync.Pool{
		New: func() any {
			start := time.Now()
			var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowCleartextPasswords=true&checkConnLiveness=true", o.User, o.Pwd, o.Address, o.Port, o.Db) //"user=" + o.User + " host=" + o.Address + " port=" + o.Port + " dbname=" + o.Db
			if o.CharSet != "" {
				dsn = dsn + "&charset=" + o.CharSet
			}
			if o.ReadTimeout > 0 {
				dsn = dsn + "&readTimeout=" + strconv.Itoa(o.ReadTimeout) + "ms"
			}
			if o.ConnTimeout > 0 {
				dsn = dsn + "&timeout=" + strconv.Itoa(o.ConnTimeout) + "ms"
			}
			if o.WriteTimeout > 0 {
				dsn = dsn + "&writeTimeout=" + strconv.Itoa(o.WriteTimeout) + "ms"
			}
			if o.Collation != "" {
				dsn = dsn + "&collation=" + o.Collation
			}
			hook := new(Hooks)
			if o.SqlLog == "yes" {
				hook.IsPrintSQLDuration = true
			}
			sql.Register("mysqlWithHooks", sqlhooks.Wrap(&mysql.MySQLDriver{}, hook))
			n, e := sqlx.Open("mysqlWithHooks", dsn)
			logger.AddMysqlConnTime(int(time.Since(start).Milliseconds()))
			if e != nil {
				logger.AddError(zap.String("dsn link fail ", o.Name+"|"+o.Address), zap.Error(e))
				panic(e.Error())
			}
			if e != nil {
				logger.AddError(zap.String("dsn ping fail ", o.Name+"|"+o.Address), zap.Error(e))
				panic(e.Error())
			}

			if o.MaxIdleConns > 0 {
				n.SetMaxIdleConns(o.MaxIdleConns)
			}
			if o.MaxOpenConns > 0 {
				n.SetMaxOpenConns(o.MaxOpenConns)
			}
			return n
		},
	}
	return c
}
