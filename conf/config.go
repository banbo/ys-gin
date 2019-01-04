package conf

import (
	beeConfig "github.com/astaxie/beego/config"
	"github.com/gin-gonic/gin"
)

var Configer *config

type config struct {
	ApiConf     ApiConfig
	DbConf      DbConfig
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

	Configer.load()
}

//加载配置到内存
func (c *config) load() {
	var err error

	//系统配置
	c.ApiConf.HttpPort = c.BeeConfiger.String("system::http_port")
	c.ApiConf.RpcPort = c.BeeConfiger.String("system::rpc_port")
	c.ApiConf.RunMode = c.BeeConfiger.String("system::run_mode")
	c.ApiConf.ParamSecret = c.BeeConfiger.String("system::param_secret")

	if c.BeeConfiger.String("system::worker_id") != "" {
		c.ApiConf.WorkerID, err = c.BeeConfiger.Int64("system::worker_id")
		if err != nil {
			panic("读取system::worker_id配置出错")
		}
	}

	//日志配置
	c.ApiConf.LogPath = c.BeeConfiger.String("log::path")
	c.ApiConf.LogLevel = c.BeeConfiger.String("log::level")

	//数据库配置
	if c.BeeConfiger.String("db::driver_name") != "" {
		c.DbConf.DriverName = c.BeeConfiger.String("db::driver_name")
		c.DbConf.Host = c.BeeConfiger.String("db::host")
		c.DbConf.Port = c.BeeConfiger.String("db::port")
		c.DbConf.User = c.BeeConfiger.String("db::user")
		c.DbConf.Password = c.BeeConfiger.String("db::password")
		c.DbConf.Database = c.BeeConfiger.String("db::database")
		c.DbConf.MaxOpen, err = c.BeeConfiger.Int("db::max_open")
		if err != nil {
			panic("读取db:max_open配置出错")
		}

		c.DbConf.MaxIdle, err = c.BeeConfiger.Int("db::max_idle")
		if err != nil {
			panic("读取db:max_idle配置出错")
		}
	}

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
