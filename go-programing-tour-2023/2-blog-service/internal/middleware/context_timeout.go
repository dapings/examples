package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func CtxTimeout(d time.Duration) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(ctx.Request.Context(), d)
		defer cancel()

		// 验证：将设置了超时的 ctx.Request.Context() 传递进去
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}
