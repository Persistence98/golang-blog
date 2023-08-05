package controller

import (
	"blogs/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Data struct {
	Article  []Article `json:"article"`
	CateList []Cate    `json:"cate_list"`
	Tag      []Tag     `json:"tag"`
	File     []Files   `json:"file"`
}

type Article struct {
	CateName  string    `json:"cate_name"`
	Title     string    `json:"title"`
	Avatar    string    `json:"avatar,omitempty"`
	ReadTime  string    `json:"read_time,omitempty"`
	Views     int       `json:"views"`
	Sort      int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cate struct {
	Id        int       `json:"id"`
	CateName  string    `json:"cate_name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tag struct {
	ID        int       `json:"id"`
	TagName   string    `json:"tag_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Files struct {
	Id        int       `json:"id"`
	FileName  string    `json:"file_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Pages struct {
	Page  int `json:"page" default:"1"`
	Limit int `json:"limit" default:"7"`
}

func Index(context *gin.Context) {
	var data Data
	var page Pages
	if err := context.ShouldBindJSON(&page); err != nil {
		if page.Page == 0 {
			page.Page = 1
		}
		if page.Limit == 0 {
			page.Limit = 7
		}
	}

	res := config.DB.Table("blog_article as article").
		Joins("join blog_cate as cate on article.cate_id = cate.id").
		Select("article.*, cate.cate_name").
		Where("article.status = ?", 1).
		Where("cate.status = ?", 1).
		Order("article.sort desc").
		Offset((page.Page - 1) * 1).
		Limit(page.Limit).
		Find(&data.Article)
	if res.Error != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  500,
			"data:":   nil,
			"message": "服务器错误！",
		})
		return
	}
	config.DB.Table("blog_cate").Where("status = ?", 1).Order("sort asc").Find(&data.CateList)
	config.DB.Table("blog_tag").Where("status = ?", 1).Order("sort asc").Find(&data.Tag)
	config.DB.Table("blog_files").Where("status = ?", 1).Order("sort asc").Find(&data.File)
	config.Redis.Get(context, "hello").Val()
	context.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功！",
		"data":    data,
	})
}
