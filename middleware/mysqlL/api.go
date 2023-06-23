package mysqlL

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/flyerxp/lib/app"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/nacos"
	yaml2 "github.com/flyerxp/lib/utils/yaml"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/qustavo/sqlhooks/v2"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

// Mysql 容器
type SqlContainer struct {
	SqlContainer cmap.ConcurrentMap[string, *MysqlClient]
	MysqlConf    cmap.ConcurrentMap[string, config2.MidMysqlConf]
}

// Mysql 客户端
type MysqlClient struct {
	Poll   *sync.Pool
	CurrDb *sqlx.DB
}

var MysqlEngine *SqlContainer

type MysqlLog struct {
}

func (m *MysqlLog) Print(v ...interface{}) {
	zapLog := make([]zap.Field, len(v))
	for i := range v {
		switch v[i].(type) {
		case error:
			zapLog[i] = zap.Error(v[i].(error))
		case string:
			zapLog[i] = zap.String("mysql driver error", v[i].(string))
		default:
			zapLog[i] = zap.Any("mysql driver error", v[i])
		}
	}
	logger.AddError(zapLog...)
}
func GetEngine(name string, ctx context.Context) (*MysqlClient, error) {
	if MysqlEngine == nil {
		MysqlEngine = new(SqlContainer)
		var confList []config2.MidMysqlConf
		MysqlEngine.MysqlConf = cmap.New[config2.MidMysqlConf]()
		MysqlEngine.SqlContainer = cmap.New[*MysqlClient]()
		conf := config2.GetConf()
		confList = conf.Mysql
		//本地文件中获取
		for _, v := range confList {
			if v.Name != "" {
				MysqlEngine.MysqlConf.Set(v.Name, v)
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
							MysqlEngine.MysqlConf.Set(v.Name, v)
						}
					} else {
						logger.AddError(zap.Error(errors.New("yaml conver error")))
					}
				}
			}
		}
		_ = app.RegisterFunc("mysql", "mysql close", func() {
			MysqlEngine.Reset()
		})
	}

	e, ok := MysqlEngine.SqlContainer.Get(name)
	if ok {
		return e, nil
	}
	o, okC := MysqlEngine.MysqlConf.Get(name)
	if okC {
		objMysql := newClient(o)
		MysqlEngine.SqlContainer.Set(name, objMysql)
		return objMysql, nil
	}
	logger.AddError(zap.Error(errors.New("no find mysql config " + name)))
	return nil, errors.New("no find mysql config " + name)
}

// https://github.com/golang-migrate/migrate/blob/master/database/mysql/README.md
func newClient(o config2.MidMysqlConf) *MysqlClient {
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
			//_ = mysql.SetLogger(&MysqlLog{})
			sql.Register("mysqlWithHooks", sqlhooks.Wrap(&mysql.MySQLDriver{}, hook))
			n, e := sqlx.Open("mysqlWithHooks", dsn)
			go func() {
				if n.Ping() != nil {
					logger.AddError(zap.Error(errors.New("dsn link fail:" + o.Address)))
				}
			}()
			logger.AddMysqlConnTime(int(time.Since(start).Milliseconds()))
			if e != nil {
				logger.AddError(zap.String("dsn link fail ", o.Name+"|"+o.Address), zap.Error(e))
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
	return &MysqlClient{c, nil}
}
func (m *MysqlClient) GetDb() *sqlx.DB {
	if m.CurrDb == nil {
		m.CurrDb = m.Poll.Get().(*sqlx.DB)
	}
	return m.CurrDb
}
func (m *MysqlClient) PutDb(a *sqlx.DB) {
	m.Poll.Put(a)
}
func (m *SqlContainer) Reset() {
	for _, v := range MysqlEngine.SqlContainer.Items() {
		_ = v.CurrDb.Close()
	}
	MysqlEngine = nil
}
