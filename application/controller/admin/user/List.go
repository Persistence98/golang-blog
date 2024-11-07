package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func List(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user/list.html", gin.H{
		"title": "User",
	})
}
