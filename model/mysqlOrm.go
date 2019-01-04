package model

import (
	"fmt"
	"log"
	"github.com/banbo/ys-gin/utils/aes"

	"github.com/banbo/ys-gin/conf"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func NewMysqlOrm() {
	var err error
	var password string
	encrypt := conf.Configer.BeeConfiger.DefaultBool("system::pswd_encrypt", false)
	if encrypt {
		pswd, err := aes.Decrypt([]byte(conf.Configer.DbConf.Password))
		if err != nil {
			log.Fatal("decrypt pswd err：", err)
		}
		password = string(pswd)
	} else {
		password = conf.Configer.DbConf.Password
	}
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		conf.Configer.DbConf.User, password, conf.Configer.DbConf.Host, conf.Configer.DbConf.Port, conf.Configer.DbConf.Database)
	Orm, err = xorm.NewEngine("mysql", conn)
	if err != nil {
		log.Fatal(err.Error())
	}

	Orm.SetMaxOpenConns(conf.Configer.DbConf.MaxOpen)
	Orm.SetMaxIdleConns(conf.Configer.DbConf.MaxIdle)

	//非生产环境显示sql
	if gin.Mode() != gin.ReleaseMode {
		Orm.ShowSQL(true)
	}
}
