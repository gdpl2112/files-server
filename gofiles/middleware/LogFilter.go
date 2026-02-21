package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LogFilter 日志中间件
func LogFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		ctx.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := ctx.Request.Method

		// 请求路由
		reqUri := ctx.Request.RequestURI

		// 状态码
		statusCode := ctx.Writer.Status()

		// 请求IP
		clientIP := ctx.ClientIP()

		// 日志格式
		fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
