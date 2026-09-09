package conf

import (
	"fmt"
	"path/filepath"
	"strconv"
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
	c.DbConf = make([]DbConfig, 0)

	// 尝试从 db 数组读取
	if dbList, err := c.BeeConfiger.DIY("db"); err == nil {
		if arr, ok := dbList.([]interface{}); ok {
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					alias := "default"
					if v, ok := m["alias"]; ok {
						alias = fmt.Sprintf("%v", v)
					}

					dbConfig := DbConfig{
						Alias:      alias,
						DriverName: fmt.Sprintf("%v", m["driver_name"]),
						Database:   fmt.Sprintf("%v", m["database"]),
						Host:       fmt.Sprintf("%v", m["host"]),
						Port:       fmt.Sprintf("%v", m["port"]),
						User:       fmt.Sprintf("%v", m["user"]),
						Password:   fmt.Sprintf("%v", m["password"]),
						Charset:    fmt.Sprintf("%v", m["charset"]),
					}

					if v, ok := m["max_open"]; ok {
						dbConfig.MaxOpen, _ = strconv.Atoi(fmt.Sprintf("%v", v))
					}
					if v, ok := m["max_idle"]; ok {
						dbConfig.MaxIdle, _ = strconv.Atoi(fmt.Sprintf("%v", v))
					}

					c.DbConf = append(c.DbConf, dbConfig)
				}
			}
			return
		}
	}
}
