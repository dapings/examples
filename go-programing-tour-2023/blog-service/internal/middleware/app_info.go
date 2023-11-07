package middleware

import (
	"github.com/gin-gonic/gin"
)

func AppInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 使用 ctx.Keys 进行元数据的存储
		ctx.Set("app_name", "blog-service")
		ctx.Set("app_version", "v1.0.0")
		ctx.Next()
	}
}
