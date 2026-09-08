package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/banbo/ys-gin/conf"
	"github.com/banbo/ys-gin/errors"
)

var Engineer *engine

type engine struct {
	gormDB map[string]*gorm.DB
}

func NewEngine() (*engine, error) {
	Engineer = &engine{
		gormDB: make(map[string]*gorm.DB),
	}

	var err error
	for _, db := range conf.Configer.DbConf {
		dbSource, err := getConnects(db)
		if err != nil {
			break
		}

		gormDB, err := gorm.Open(db.DriverName, dbSource)
		if err != nil {
			break
		}

		//连接是否可用
		if err = gormDB.DB().Ping(); err != nil {
			break
		}

		//最大打开连接数
		gormDB.DB().SetMaxOpenConns(db.MaxOpen)
		//连接池的空闲数大小
		gormDB.DB().SetMaxIdleConns(db.MaxIdle)

		//结构体命名与数据库一致
		gormDB.SingularTable(true)

		Engineer.gormDB[db.Alias] = gormDB
	}

	return Engineer, err
}

//获取数据库引擎
func (e *engine) Get(mi ModelInterface) (*gorm.DB, error) {
	alias := mi.DatabaseAlias()
	db, ok := e.gormDB[alias]
	if !ok {
		return nil, errors.NewNormal("数据库引擎：" + alias + "不存在")
	}

	return db, nil
}

// 获取数据库连接信息
func getConnects(db conf.DbConfig) (string, error) {
	var err error
	var dataSource string

	switch db.DriverName {
	case "mysql":
		dataSource = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			db.User, db.Password, db.Host, db.Port, db.Database, db.Charset)
	case "postgres":
		dataSource = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password, db.Database)
	case "sqlite3":
		dataSource = db.Database
	default:
		err = errors.New("不支持的数据库类型：" + db.DriverName)
	}

	return dataSource, err
}
