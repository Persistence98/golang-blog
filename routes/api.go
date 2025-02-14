package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-blog/application/controller/admin/article"
	"golang-blog/application/controller/admin/index"
	Login "golang-blog/application/controller/admin/login"
	"golang-blog/application/controller/admin/user"
	"golang-blog/application/controller/api/v1"
	"golang-blog/middleware"
)

func init() {
	routes := gin.Default()

	routes.Use(middleware.HeadersMiddleware()) //请求头中间件
	initAdminRoutes(routes)                    //初始化后台路由

	routes.Use(middleware.RateLimiterMiddleware()) //令牌桶中间件
	initApiRoutes(routes)                          //初始化接口路由

	err := routes.Run(":1040")
	if err != nil {
		fmt.Printf(`路由初始化失败！`)
	}
}

/*接口路由*/
func initApiRoutes(engine *gin.Engine) {
	api := engine.Group("api")
	api.Use(middleware.JwtMiddleware()) //jwt鉴权
	{
		route := api.Group("v1")
		route.POST("/index", v1.Index)
	}
}

/*后台路由*/
func initAdminRoutes(engine *gin.Engine) {
	engine.LoadHTMLGlob("templates/**/*")
	engine.Static("/static", "public/assets")
	api := engine.Group("blog_admin")
	{ /* 登陆模块 */
		route := api.Group("login")
		route.GET("", Login.Index)
		route.POST("login", Login.Login)
		route.GET("getCodeImg", Login.GetCodeImg)
		route.POST("reloadCaptcha", Login.ReloadCaptcha)
	}
	api.Use(middleware.AdminAuthMiddle())
	{ /* 首页框架 */
		route := api.Group("home")
		route.GET("index", index.Index)
		route.GET("console", index.Console)
	}
	{ /* 用户管理 */
		route := api.Group("users")
		route.GET("list", user.List)
		route.GET("list_data", user.ListData)
		route.GET("DownLoadExcel", user.DownLoadExcel)
	}
	{ /* 内容管理 */
		route := api.Group("article")
		route.GET("cate_list", article.CateList)
	}
}
