package middleware

import (
	"github.com/gin-gonic/gin"
)

func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
