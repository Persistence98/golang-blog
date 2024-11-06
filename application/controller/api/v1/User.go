package v1

import (
	"github.com/gin-gonic/gin"
	"golang-blog/config"
	"net/http"
	"time"
)

type Users struct {
	Id        int       `json:"id" binding:"required"`
	Nickname  string    `json:"nickname"`
	UserName  string    `json:"user_name"`
	Password  string    `json:"password"`
	Sex       int       `json:"sex"`
	Age       int       `json:"age"`
	Address   string    `json:"address"`
	Status    int32     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DownloadUserExcel(ctx *gin.Context) {
	var user Users
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "参数不全！",
		})
		return
	}
	config.DB.Table("users").Where("id = ?", user.Id).First(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    user,
		"message": "请求成功！",
	})
}
