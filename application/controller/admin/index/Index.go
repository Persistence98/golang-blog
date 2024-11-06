package index

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}

func Console(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "console.html", gin.H{})
}
