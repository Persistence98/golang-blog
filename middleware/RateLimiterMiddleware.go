package middleware

// 注册令牌桶 设置每秒生成令牌数、令牌桶总容量
//var limiter = ratelimit.NewBucketWithRate(400, 1800)
//
//func RateLimiterMiddleware() gin.HandlerFunc {
//	return func(context *gin.Context) {
//		if limiter.TakeAvailable(1) < 1 {
//			context.JSON(http.StatusTooManyRequests, gin.H{
//				"status":  429,
//				"message": "压力山大啊,请休息一会再看看吧！",
//			})
//			context.Abort()
//			return
//		}
//		context.Next()
//	}
//}
