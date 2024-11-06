package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Pages struct {
	Page  int `json:"page" default:"1"`
	Limit int `json:"limit" default:"7"`
}

func Index(context *gin.Context) {
	//var page Pages

	context.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "获取成功！",
		"data":    []int{},
	})
}

func Index2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "获取成功！",
		"data":    []int{},
	})
}
