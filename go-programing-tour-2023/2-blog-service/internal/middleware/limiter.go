package middleware

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := l.Key(ctx)
		if bucket, ok := l.GetBucket(key); ok {
			// 占用存储桶中立即可用的令牌的数量
			count := bucket.TakeAvailable(1)
			if count == 0 {
				resp := app.NewRes(ctx)
				resp.ToErrRes(errcode.TooManyRequests)
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}
