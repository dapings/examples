package log

import (
	"context"

	"go-micro.dev/v4/server"
	"go-micro.dev/v4/util/log"
	"go.uber.org/zap"
)

// micro 中间件：对请求进行一层封装，在接收到gRPC请求时，输出请求的具体参数

func LogWrapper(logger *zap.Logger) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			log.Info("receive request",
				zap.String("method", req.Method()),
				zap.String("Service", req.Service()),
				zap.Reflect("request param:", req.Body()),
			)
			return fn(ctx, req, rsp)
		}
	}
}
