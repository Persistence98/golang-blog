// 在 middleware 中设置请求头
package middleware

import "github.com/gin-gonic/gin"

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置跨域请求头
		c.Header("Access-Control-Allow-Origin", "*") // 允许所有域的跨域请求
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		// 如果你需要设置授权 token 等
		// c.Header("Authorization", "Bearer your_token")

		// 调用下一个中间件/处理函数
		c.Next()
	}
}
