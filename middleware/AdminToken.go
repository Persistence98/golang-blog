package middleware

import (
	"github.com/gin-gonic/gin"
	"golang-blog/config"
	"net/http"
)

func AdminAuthMiddle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			ctx.Redirect(http.StatusFound, "/blog_admin/login")
			ctx.Abort()
		}

		id, err := config.Redis.HGet(ctx, "token", token).Result()
		if err != nil {
			ctx.Redirect(http.StatusFound, "/blog_admin/login")
			ctx.Abort()
		}
		ctx.Set("user_id", id)

		ctx.Next()
	}
}
