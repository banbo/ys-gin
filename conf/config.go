package conf

import (
	"fmt"
	"strings"

	beeConfig "github.com/astaxie/beego/config"
	"github.com/gin-gonic/gin"
)

var Configer *config

type config struct {
	ApiConf     ApiConfig
	DbConf      []DbConfig
	RedisConf   RedisConfig
	BeeConfiger beeConfig.Configer
}

func NewConfiger(filename string) {
	Configer = new(config)

	var err error
	Configer.BeeConfiger, err = beeConfig.NewConfig("ini", filename)
	if err != nil {
		panic("读取配置文件出错")
	}

	//读取配置
	Configer.load()

	Configer.loadDbs()
}

//加载配置到内存
func (c *config) load() {
	var err error

	//系统配置
	c.ApiConf.HttpPort = c.BeeConfiger.String("system::http_port")
	c.ApiConf.RpcPort = c.BeeConfiger.String("system::rpc_port")
	c.ApiConf.RunMode = c.BeeConfiger.String("system::run_mode")
	c.ApiConf.ParamSecret = c.BeeConfiger.String("system::param_secret")
	c.ApiConf.Dbs = c.BeeConfiger.String("system::dbs")

	if c.BeeConfiger.String("system::worker_id") != "" {
		c.ApiConf.WorkerID, err = c.BeeConfiger.Int64("system::worker_id")
		if err != nil {
			panic("读取system::worker_id配置出错")
		}
	}

	//日志配置
	c.ApiConf.LogPath = c.BeeConfiger.String("log::path")
	c.ApiConf.LogLevel = c.BeeConfiger.String("log::level")

	//redis配置
	if c.BeeConfiger.String("redis::host") != "" {
		c.RedisConf.Host = c.BeeConfiger.String("redis::host")
		c.RedisConf.Port = c.BeeConfiger.String("redis::port")
		c.RedisConf.Password = c.BeeConfiger.String("redis::password")
		c.RedisConf.DB, err = c.BeeConfiger.Int("redis::db")
		if err != nil {
			panic("读取redis::db配置出错")
		}
	}

	//判断配置
	if c.ApiConf.RunMode != gin.DebugMode && c.ApiConf.RunMode != gin.TestMode && c.ApiConf.RunMode != gin.ReleaseMode {
		panic("run_mode配置错误")
	}
	if c.ApiConf.LogLevel != "debug" && c.ApiConf.LogLevel != "info" && c.ApiConf.LogLevel != "error" {
		panic("log_level配置错误")
	}
}

//读取数据库配置
func (c *config) loadDbs() {
	// 数据库配置
	dbs := strings.Split(c.ApiConf.Dbs, ",")
	c.DbConf = make([]DbConfig, 0, len(dbs))
	if len(dbs) > 0 && len(dbs[0]) > 0 {
		for _, db := range dbs {
			prefix := "db-" + db

			//获取最大连接数，如果配置了
			var maxOpen, maxIdle int
			var err error
			if c.BeeConfiger.String(prefix+"::max_open") != "" {
				maxOpen, err = c.BeeConfiger.Int(prefix + "::max_open")
				if err != nil {
					panic("读取db::max_open配置出错")
				}

				maxIdle, err = c.BeeConfiger.Int(prefix + "::max_idle")
				if err != nil {
					panic("读取db::max_idle配置出错")
				}
			}

			dbConfig := DbConfig{
				Alias:      db,
				DriverName: c.BeeConfiger.String(prefix + "::driver_name"),
				Database:   c.BeeConfiger.String(prefix + "::database"),
				Host:       c.BeeConfiger.String(prefix + "::host"),
				Port:       c.BeeConfiger.String(prefix + "::port"),
				User:       c.BeeConfiger.String(prefix + "::user"),
				Password:   c.BeeConfiger.String(prefix + "::password"),
				Charset:    c.BeeConfiger.String(prefix + "::charset"),
				MaxOpen:    maxOpen,
				MaxIdle:    maxIdle,
			}

			//slaves
			for i := 0; i < 10; i++ {
				prefixSlave := prefix + "-slave-" + fmt.Sprint(i)
				slaveHost := c.BeeConfiger.String(prefixSlave + "::host")
				if slaveHost == "" {
					break
				}

				dbConfig.Slaves = append(dbConfig.Slaves, struct {
					Host     string
					Port     string
					User     string
					Password string
				}{Host: slaveHost, Port: c.BeeConfiger.String(prefixSlave + "::port"), User: c.BeeConfiger.String(prefixSlave + "::user"), Password: c.BeeConfiger.String(prefixSlave + "::password")})
			}

			c.DbConf = append(c.DbConf, dbConfig)
		}
	}
}
