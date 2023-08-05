package config

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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
	SetConnectMysqlMaxPool(connect)
}

func SetConnectMysqlMaxPool(g *gorm.DB) {
	mysql, err := g.DB()
	if err != nil {
		panic(err)
	}

	//设置MySQL连接池最大连接数量
	mysql.SetMaxOpenConns(30)
	//设置MySQL连接池最大空闲数
	mysql.SetConnMaxIdleTime(7)
	// 设置连接的最大可复用时间
	mysql.SetConnMaxLifetime(time.Hour)
}
