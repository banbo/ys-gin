package cache

import (
	"github.com/banbo/ys-gin/conf"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func NewRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Configer.RedisConf.Host + ":" + conf.Configer.RedisConf.Port,
		Password: conf.Configer.RedisConf.Password,
		DB:       conf.Configer.RedisConf.DB,
	})
}
