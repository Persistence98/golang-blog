package config

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("打开配置文件错误！")
	}

	dns := viper.GetString("database.mysql")

	connect, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败！")
	}
	DB = connect
}
