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

		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}
