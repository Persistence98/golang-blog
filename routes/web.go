package routes

import (
	"blogs/application/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WebRoutes(r *gin.Engine) {
	//r.Use(middleware.RateLimiterMiddleware()) //令牌桶
	//r.Use(middleware.LoggingMiddleware())     //用户来源

	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "初始化成功！"})
	})
	api := r.Group("api")
	{
		v1 := api.Group("v1")
		v1.POST("/index", controller.Index)
		v1.POST("/download_user_excel", controller.DownloadUserExcel)
	}
}
