package Login

import (
	"github.com/dchest/captcha"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	Models "golang-blog/application/models"
	"golang-blog/config"
	"net/http"
	"time"
)

func Index(ctx *gin.Context) {
	captchaID := NewCaptcha(ctx) // 生成新的验证码ID

	token, err := ctx.Cookie("token")
	if err == nil && token != "" {
		config.Redis.HDel(ctx, "token", token)
	}
	ctx.HTML(200, "login.html", gin.H{
		"captchaID": captchaID,
	})
}

func NewCaptcha(ctx *gin.Context) string {
	return captcha.New()
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

/**
*	刷新二维码
**/
func ReloadCaptcha(ctx *gin.Context) {
	capId := NewCaptcha(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   capId,
		"msg":    "获取成功！",
	})
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
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"status": http.StatusNotAcceptable,
			"msg":    "验证码错误",
		})
		return
	}

	var user Models.User
	if err := config.DB.Where("user_name = ?", log.UserName).Where("status = ?", 1).Select("id,nickname,user_name,password,created_at").First(&user).Error; err != nil {
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

	jwtString, err := generateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    err.Error(),
		})
		return
	}
	config.Redis.HSet(ctx, "token", jwtString, user.ID)
	ctx.SetCookie("token", jwtString, 86400, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "登陆成功！",
	})
	return
}

/* 生成jwt */
func generateJWT(user Models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	Claims := token.Claims.(jwt.MapClaims)

	Claims["sub"] = user.ID
	//Claims["nickname"] = user.Nickname
	//Claims["user_name"] = user.UserName

	Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return "", err
	}
	secretKey := []byte(viper.GetString("middleware.jwt-secret-key"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
