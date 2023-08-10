package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		viper.SetConfigFile("config.yaml")
		err := viper.ReadInConfig()
		if err != nil {
			panic("打开配置文件失败！")
		}
		secretKey := viper.GetString("middleware.jwt-secret-key")

		AccessUserToken := ctx.GetHeader("Access-User-Token") //获取请求头

		if AccessUserToken == "" { //判断是否为空或者没有请求头
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "请重新登陆！",
			})
			ctx.Abort()
			return
		}

		token, errors := jwt.Parse(AccessUserToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if errors != nil || !token.Valid {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "请重新登陆！",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
