package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FileValidationFilter 文件验证中间件
func FileValidationFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 只对上传请求进行验证
		if ctx.Request.Method == http.MethodPost && (ctx.Request.URL.Path == "/upload" || ctx.Request.URL.Path == "/user/upload") {
			file, err := ctx.FormFile("file")
			if err == nil {
				// 检查文件大小，这里设置一个默认限制，例如100MB
				maxSize := int64(100 * 1024 * 1024)
				if file.Size > maxSize {
					ctx.JSON(http.StatusBadRequest, gin.H{"message": "文件大小超过限制"})
					ctx.Abort()
					return
				}
			}
		}

		ctx.Next()
	}
}
