package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/banbo/ys-gin/conf"
	"github.com/banbo/ys-gin/errors"
)

var Engineer *engine

type engine struct {
	xormEngine map[string]*xorm.Engine
}

func NewEngine() (*engine, error) {
	Engineer = &engine{
		xormEngine: make(map[string]*xorm.Engine),
	}

	var err error
	var dataSourceName string
	var xormEngine *xorm.Engine
	for _, db := range conf.Configer.DbConf {
		dataSourceName, err = getConnect(db)
		if err != nil {
			break
		}

		xormEngine, err = xorm.NewEngine(db.DriverName, dataSourceName)
		if err != nil {
			break
		}

		//连接是否可用
		if err = xormEngine.Ping(); err != nil {
			break
		}

		//最大打开连接数
		xormEngine.SetMaxOpenConns(db.MaxOpen)
		//连接池的空闲数大小
		xormEngine.SetMaxIdleConns(db.MaxIdle)

		//结构体命名与数据库一致
		xormEngine.SetMapper(core.NewCacheMapper(new(core.SameMapper)))

		Engineer.xormEngine[db.Alias] = xormEngine
	}

	return Engineer, err
}

//获取数据库引擎
func (e *engine) Get(mi ModelInterface) (*xorm.Engine, error) {
	alias := mi.DatabaseAlias()
	engine, ok := e.xormEngine[alias]
	if !ok {
		return nil, errors.NewNormal("数据库引擎：" + alias + "不存在")
	}

	return engine, nil
}

// 获取数据库连接信息
func getConnect(db conf.DbConfig) (string, error) {
	var err error
	var dataSourceName string

	switch db.DriverName {
	case "mysql":
		dataSourceName = db.User + ":" + db.Password
		dataSourceName += "@(" + db.Host + ":" + db.Port + ")/"
		dataSourceName += db.Database + "?charset=utf8&loc=Asia%2FShanghai&multiStatements=true"
	case "postgres":
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password, db.Database)
	case "sqlite3":
		dataSourceName = db.Database
	default:
		err = errors.New("不支持的数据库类型：" + db.DriverName)
	}

	return dataSourceName, err
}
