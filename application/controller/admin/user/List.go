package user

import (
	"github.com/gin-gonic/gin"
	Models "golang-blog/application/models"
	"golang-blog/config"
	"net/http"
)

func List(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user/list.html", gin.H{
		"title": "User",
	})
}

type Pages struct {
	Page    int    `json:"page,omitempty" form:"page"`
	Limit   int    `json:"limit,omitempty" form:"limit"`
	Keyword string `json:"keyword,omitempty" form:"keyword"`
	State   int    `json:"state" form:"state"`
}

func ListData(ctx *gin.Context) {
	var pages Pages
	if err := ctx.ShouldBind(&pages); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": err.Error()})
		return
	}
	pageNum := (pages.Page - 1) * pages.Limit

	query := config.DB

	if pages.Keyword != "" {
		query = query.Where("user_name like ? or address like ? or nickname like ?", pages.Keyword+"%", "%"+pages.Keyword+"%", pages.Keyword+"%")
	}
	if pages.State != 0 {
		query = query.Where("status = ?", pages.State)
	}
	countSql := query
	var data []Models.User
	if err := query.Limit(pages.Limit).Offset(pageNum).Where("status != ?", 3).Select("id,nickname,user_name,sex,age,address,status,created_at,updated_at").Order("id desc").Find(&data).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    err.Error(),
		})
		return
	}
	var count int64
	if err := countSql.Model(Models.User{}).Where("status != ?", 3).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    err.Error(),
		})
		return
	}
	//config.Redis.Del(ctx, "user_count")
	//count, err := config.Redis.Get(ctx, "user_count").Int64()
	//if err == redis.Nil {
	//
	//	config.Redis.Set(ctx, "user_count", count, 30*time.Minute)
	//}

	sexMap := map[string]string{
		"1": "男",
		"2": "女",
	}
	statusMap := map[string]string{
		"1": "正常",
		"2": "锁定",
		"3": "删除",
	}
	for index, item := range data {
		data[index].Sex = sexMap[item.Sex]
		data[index].Status = statusMap[item.Status]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "请求成功！",
		"count": count,
		"data":  data,
	})
	return
}
