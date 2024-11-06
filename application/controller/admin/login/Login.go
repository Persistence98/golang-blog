package Login

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	Models "golang-blog/application/models"
	"golang-blog/config"
	"net/http"
)

func Index(ctx *gin.Context) {
	captchaID := captcha.New() // 生成新的验证码ID
	ctx.HTML(200, "login.html", gin.H{
		"title":     "标题",
		"captchaID": captchaID,
	})
}

// GetCodeImg 生成验证码
func GetCodeImg(ctx *gin.Context) {
	captchaID := ctx.DefaultQuery("captchaID", "")
	if captchaID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"error":  "验证码生成失败",
		})
		return
	}
	err := captcha.WriteImage(ctx.Writer, captchaID, 240, 80)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": err,
			"error":  "验证码生成失败",
		})
		return
	}
}

type LoginRequest struct {
	UserName     string `json:"user_name" binding:"required,max=12"`
	UserPassword string `json:"password" binding:"required,max=50"`
	Code         string `json:"code" binding:"required,max=6"`
	CaptchaID    string `json:"captchaID" binding:"required"`
}

func Login(ctx *gin.Context) {
	var log LoginRequest
	if err := ctx.ShouldBindJSON(&log); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    err.Error(),
		})
		return
	}

	if !captcha.VerifyString(log.CaptchaID, log.Code) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "验证码错误",
		})
		return
	}

	var user Models.User
	if err := config.DB.Where("user_name = ?", log.UserName).Where("status = ?", 1).Select("id,nickname,user_name,password").First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "账号或密码错误！",
		})
		return
	}

	if user.Password != log.UserPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "账号或密码错误！",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})
	return
}
