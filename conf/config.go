package conf

import (
	"path/filepath"
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
	configType  string
}

// 获取配置key，根据类型自动选择分隔符
func (c *config) key(section, key string) string {
	if c.configType == "yaml" || c.configType == "json" {
		return section + "." + key
	}
	return section + "::" + key
}

func NewConfiger(filename string) {
	Configer = new(config)

	// 根据文件后缀选择配置类型
	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".yaml", ".yml":
		Configer.configType = "yaml"
	case ".json":
		Configer.configType = "json"
	case ".xml":
		Configer.configType = "xml"
	default:
		Configer.configType = "ini"
	}

	var err error
	if Configer.configType == "yaml" {
		// yaml 使用自定义适配器（gopkg.in/yaml.v3）
		Configer.BeeConfiger, err = NewYamlConfiger(filename)
	} else {
		Configer.BeeConfiger, err = beeConfig.NewConfig(Configer.configType, filename)
	}
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
	c.ApiConf.HttpPort = c.BeeConfiger.String(c.key("system", "http_port"))
	c.ApiConf.RpcPort = c.BeeConfiger.String(c.key("system", "rpc_port"))
	c.ApiConf.RunMode = c.BeeConfiger.String(c.key("system", "run_mode"))
	c.ApiConf.ParamSecret = c.BeeConfiger.String(c.key("system", "param_secret"))
	c.ApiConf.Dbs = c.BeeConfiger.String(c.key("system", "dbs"))

	if c.BeeConfiger.String(c.key("system", "worker_id")) != "" {
		c.ApiConf.WorkerID, err = c.BeeConfiger.Int64(c.key("system", "worker_id"))
		if err != nil {
			panic("读取system::worker_id配置出错")
		}
	}

	//日志配置
	c.ApiConf.LogPath = c.BeeConfiger.String(c.key("log", "path"))
	c.ApiConf.LogLevel = c.BeeConfiger.String(c.key("log", "level"))

	//redis配置
	if c.BeeConfiger.String(c.key("redis", "host")) != "" {
		c.RedisConf.Host = c.BeeConfiger.String(c.key("redis", "host"))
		c.RedisConf.Port = c.BeeConfiger.String(c.key("redis", "port"))
		c.RedisConf.Password = c.BeeConfiger.String(c.key("redis", "password"))
		c.RedisConf.DB, err = c.BeeConfiger.Int(c.key("redis", "db"))
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
			if c.BeeConfiger.String(c.key(prefix, "max_open")) != "" {
				maxOpen, err = c.BeeConfiger.Int(c.key(prefix, "max_open"))
				if err != nil {
					panic("读取db::max_open配置出错")
				}

				maxIdle, err = c.BeeConfiger.Int(c.key(prefix, "max_idle"))
				if err != nil {
					panic("读取db::max_idle配置出错")
				}
			}

			dbConfig := DbConfig{
				Alias:      db,
				DriverName: c.BeeConfiger.String(c.key(prefix, "driver_name")),
				Database:   c.BeeConfiger.String(c.key(prefix, "database")),
				Host:       c.BeeConfiger.String(c.key(prefix, "host")),
				Port:       c.BeeConfiger.String(c.key(prefix, "port")),
				User:       c.BeeConfiger.String(c.key(prefix, "user")),
				Password:   c.BeeConfiger.String(c.key(prefix, "password")),
				Charset:    c.BeeConfiger.String(c.key(prefix, "charset")),
				MaxOpen:    maxOpen,
				MaxIdle:    maxIdle,
			}

			c.DbConf = append(c.DbConf, dbConfig)
		}
	}
}
