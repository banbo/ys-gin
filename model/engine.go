package model

import (
	"fmt"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

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
		dialector, err := getDialector(db)
		if err != nil {
			break
		}

		//结构体命名与数据库一致
		gormDB, err := gorm.Open(dialector, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			break
		}

		sqlDB, err := gormDB.DB()
		if err != nil {
			break
		}

		//连接是否可用
		if err = sqlDB.Ping(); err != nil {
			break
		}

		//最大打开连接数
		sqlDB.SetMaxOpenConns(db.MaxOpen)
		//连接池的空闲数大小
		sqlDB.SetMaxIdleConns(db.MaxIdle)

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

// 获取数据库驱动
func getDialector(db conf.DbConfig) (gorm.Dialector, error) {
	var err error
	var dataSource string

	switch db.DriverName {
	case "mysql":
		loc := db.Loc
		if loc == "" {
			loc = "Local"
		}
		dataSource = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
			db.User, db.Password, db.Host, db.Port, db.Database, db.Charset, url.QueryEscape(loc))
		return mysql.Open(dataSource), nil
	case "postgres":
		dataSource = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password, db.Database)
		return postgres.Open(dataSource), nil
	case "sqlite3":
		dataSource = db.Database
		return sqlite.Open(dataSource), nil
	default:
		err = errors.New("不支持的数据库类型：" + db.DriverName)
	}

	return nil, err
}
