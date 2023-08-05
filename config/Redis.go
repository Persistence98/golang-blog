package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Redis *redis.Client

func init() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("初始化redis配置失败！")
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.Addr"),
		Password:     viper.GetString("redis.Pwd"),
		PoolSize:     30,
		MinIdleConns: 7,
	})
}
