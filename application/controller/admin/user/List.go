package user

import (
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	Models "golang-blog/application/models"
	"golang-blog/config"
	"net/http"
	"strconv"
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
		query = query.Where("user_name like ? or address like ? or nickname like ?", pages.Keyword+"%", pages.Keyword+"%", pages.Keyword+"%")
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

// DownLoadExcel 下载用户数据
func DownLoadExcel(ctx *gin.Context) {
	//创建一个新excel
	excel := excelize.NewFile()
	sheetName := "sheet1" //文件名称
	defer func() {        //关闭文档
		if err := excel.Close(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"msg":    err.Error(),
			})
			return
		}
	}()

	streamWriter, err := excel.NewStreamWriter(sheetName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status1": http.StatusInternalServerError,
			"msg":     err.Error(),
		})
		return
	}
	defer func(streamWriter *excelize.StreamWriter) { //处理关闭错误
		err := streamWriter.Flush()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status1": http.StatusInternalServerError,
				"msg":     err.Error(),
			})
		}
	}(streamWriter)

	header := []interface{}{"ID", "用户名", "账号", "性别", "住址", "状态", "创建时间"}
	err = streamWriter.SetRow("A1", header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status2": http.StatusInternalServerError,
			"msg":     err.Error(),
		})
		return
	}

	// 分批查询数据
	limit := 200000 // 每次查询 20 万条
	rowIndex := 2   // Excel 从 2 开始填充数据（1 是表头）
	var lastID int
	maxFor := 0
	for {
		var users []Models.User
		result := config.DB.Where("status in ?", []int{1, 2}).
			Where("id > ?", lastID).
			Select("id,nickname,user_name,sex,address,status,created_at").
			Limit(limit).
			Find(&users)

		if result.Error != nil || len(users) == 0 || maxFor == 5 {
			break
		}

		for _, item := range users {
			row := []interface{}{item.ID, item.Nickname, item.UserName, item.Sex, item.Address, item.Status, item.CreatedAt}
			cell := "A" + strconv.Itoa(rowIndex)
			err := streamWriter.SetRow(cell, row)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status4": http.StatusInternalServerError, "msg": err.Error()})
				return
			}
			rowIndex++
		}
		// 更新 lastID 为当前查询结果中的最大 id
		lastID = users[len(users)-1].ID
		maxFor++
	}

	// 保存 Excel 到内存并返回给客户端
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=users.xlsx")
	ctx.Header("File-Transfer-Encoding", "binary")
	if err := excel.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
	}

}
