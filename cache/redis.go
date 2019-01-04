package cache

import (
	"github.com/go-redis/redis"
	"log"
	"github.com/banbo/ys-gin/conf"
	"github.com/banbo/ys-gin/utils/aes"
)

var RedisClient *redis.Client

func NewRedisClient() {
	var password string
	encrypt := conf.Configer.BeeConfiger.DefaultBool("system::pswd_encrypt", false)
	if encrypt {
		pswd, err := aes.Decrypt([]byte(conf.Configer.RedisConf.Password))
		if err != nil {
			log.Fatal("decrypt pswd err：", err)
		}
		password = string(pswd)
	} else {
		password = conf.Configer.RedisConf.Password
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Configer.RedisConf.Host + ":" + conf.Configer.RedisConf.Port,
		Password: password,
		DB:       conf.Configer.RedisConf.DB,
	})
}
