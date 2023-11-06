package middleware

import (
	"bytes"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	// 双写：buffer write, resp write
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bw := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		// 替换原有
		ctx.Writer = bw
		beginTime := time.Now().Unix()
		ctx.Next()
		endTime := time.Now().Unix()

		logFields := logger.Fields{
			"request":  ctx.Request.PostForm.Encode(),
			"response": bw.body.String(),
		}
		logContentFormat := "access log: method %s, status_code %d, " +
			"begin_time %d, end_time %d"
		global.Logger.WithFields(logFields).Infof(ctx, logContentFormat,
			ctx.Request.Method,
			bw.Status(),
			beginTime,
			endTime,
		)
	}
}
