package routes

import (
	"blogs/application/controller"
	"blogs/middleware"
	"github.com/gin-gonic/gin"
)

func WebRoutes(r *gin.Engine) {
	r.Use(middleware.RateLimiterMiddleware()) //令牌桶
	//r.Use(middleware.LoggingMiddleware())     //用户来源
	api := r.Group("api")
	{
		v1 := api.Group("v1")
		//v1.Use(middleware.JwtMiddleware()) //jwt鉴权
		v1.POST("/index", controller.Index)
		v1.POST("/download_user_excel", controller.DownloadUserExcel)
	}
}
