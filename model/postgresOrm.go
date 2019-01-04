package model

import (
	"fmt"
	"log"
	"github.com/banbo/ys-gin/utils/aes"

	"github.com/banbo/ys-gin/conf"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var Orm *xorm.Engine

func NewPostgresOrm() {
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
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Configer.DbConf.Host, conf.Configer.DbConf.Port, conf.Configer.DbConf.User, password, conf.Configer.DbConf.Database)
	Orm, err = xorm.NewEngine("postgres", psqlInfo)

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
