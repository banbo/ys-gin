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
	xormEngineGroup map[string]*xorm.EngineGroup
}

func NewEngine() (*engine, error) {
	Engineer = &engine{
		xormEngineGroup: make(map[string]*xorm.EngineGroup),
	}

	var err error
	for _, db := range conf.Configer.DbConf {
		dataSourceNames, err := getConnects(db)
		if err != nil {
			break
		}

		xormEngineGroup, err := xorm.NewEngineGroup(db.DriverName, dataSourceNames)
		if err != nil {
			break
		}

		//连接是否可用
		if err = xormEngineGroup.Ping(); err != nil {
			break
		}

		//最大打开连接数
		xormEngineGroup.SetMaxOpenConns(db.MaxOpen)
		//连接池的空闲数大小
		xormEngineGroup.SetMaxIdleConns(db.MaxIdle)

		//结构体命名与数据库一致
		xormEngineGroup.SetMapper(core.NewCacheMapper(new(core.SameMapper)))

		Engineer.xormEngineGroup[db.Alias] = xormEngineGroup
	}

	return Engineer, err
}

//获取数据库引擎
func (e *engine) Get(mi ModelInterface) (*xorm.EngineGroup, error) {
	alias := mi.DatabaseAlias()
	engineGroup, ok := e.xormEngineGroup[alias]
	if !ok {
		return nil, errors.NewNormal("数据库引擎：" + alias + "不存在")
	}

	return engineGroup, nil
}

// 获取数据库连接信息
func getConnects(db conf.DbConfig) ([]string, error) {
	var err error
	var dataSourceNames []string

	switch db.DriverName {
	case "mysql":
		dataSourceName := db.User + ":" + db.Password
		dataSourceName += "@(" + db.Host + ":" + db.Port + ")/"
		dataSourceName += db.Database + "?charset=" + db.Charset

		dataSourceNames = append(dataSourceNames, dataSourceName)

		for _, v := range db.Slaves {
			dataSourceName := v.User + ":" + v.Password
			dataSourceName += "@(" + v.Host + ":" + v.Port + ")/"
			dataSourceName += db.Database + "?charset=" + db.Charset

			dataSourceNames = append(dataSourceNames, dataSourceName)
		}

	case "postgres":
		dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password, db.Database)
		dataSourceNames = append(dataSourceNames, dataSourceName)

		for _, v := range db.Slaves {
			dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				v.Host, v.Port, v.User, v.Password, db.Database)

			dataSourceNames = append(dataSourceNames, dataSourceName)
		}
	case "sqlite3":
		dataSourceName := db.Database
		dataSourceNames = append(dataSourceNames, dataSourceName)
	default:
		err = errors.New("不支持的数据库类型：" + db.DriverName)
	}

	return dataSourceNames, err
}
