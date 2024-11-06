package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-blog/application/controller/api/v1"
	"golang-blog/middleware"
)

func init() {
	routes := gin.Default()

	initAdminRoutes(routes)                        //初始化后台路由
	routes.Use(middleware.RateLimiterMiddleware()) //令牌桶中间件
	initApiRoutes(routes)                          //初始化接口路由

	err := routes.Run(":1040")
	if err != nil {
		fmt.Printf(`路由初始化失败！`)
	}
}

func initAdminRoutes(engine *gin.Engine) {
	api := engine.Group("admin")
	{
		route := api.Group("article")
		route.POST("/index", v1.Index2)
	}
}

func initApiRoutes(engine *gin.Engine) {
	api := engine.Group("api")
	api.Use(middleware.JwtMiddleware()) //jwt鉴权
	{
		route := api.Group("v1")
		route.POST("/index", v1.Index)
	}
}
