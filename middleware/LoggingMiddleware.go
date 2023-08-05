package middleware

import (
	"blogs/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Log struct {
	Ip            string
	Type          string
	ContinentCode string
	ContinentName string
	CountryCode   string
	CountryName   string
	RegionCode    string
	RegionName    string
	City          string
	Zip           string
	Latitude      float64
	Longitude     float64
	Code          string
	Public        string
	CreatedAt     time.Time
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		var logs Log
		var channels = make(chan Log, 1)
		ip := context.ClientIP()
		if ip == "" || ip == "127.0.0.1" {
			context.Next()
			return
		}
		go func() {
			// 调用 ipapi 接口返回数据
			viper.SetConfigFile("config.yaml")
			AccessKey := viper.GetString("middleware.IpApi-Access-Key")
			IpApiURL := "http://api.ipapi.com/" + url.PathEscape(ip) + "?access_key" + AccessKey + "=&language=zh"
			resp, err := http.Get(IpApiURL)
			if err != nil {
				context.Next()
				return
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			res, err := io.ReadAll(resp.Body)
			if err != nil {
				context.Next()
				return
			}
			err = json.Unmarshal(res, &logs)
			if err != nil {
				context.Next()
				return
			}
			//执行完后向管道发送数据
			channels <- logs
		}()
		go func() {
			logs := <-channels
			logs.Public = context.Request.UserAgent()
			_ = config.DB.Table("blog_log").Create(&logs).Error
		}()
		context.Next()
	}
}
