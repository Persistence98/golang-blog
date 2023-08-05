package main

import (
	"blogs/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//注册路由
	routes.WebRoutes(r)
	err := r.Run(":1040")
	if err != nil {
		panic("监听端口失败！")
	}
}
