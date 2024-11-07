package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CateList(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "article/list.html", gin.H{
		"title": "Article",
	})
}
