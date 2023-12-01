package middleware

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/errcode"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/meta"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerTracing(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// 读取RPC方法传入的上下文信息，解析出metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	// 从给定的载体中解码出SpanContext实例，并创建和设置本次Span的标签信息
	parentSpanCtx, _ := global.Tracer.Extract(opentracing.TextMap, meta.MetadataTextMap{MD: md})
	spanOpts := []opentracing.StartSpanOption{
		opentracing.Tag{
			Key:   string(ext.Component),
			Value: global.GRPCSpanTagVal,
		},
		ext.SpanKindRPCServer,
		ext.RPCServerOption(parentSpanCtx),
	}
	span := global.Tracer.StartSpan(info.FullMethod, spanOpts...)
	defer span.Finish()

	// 根据当前Span返回一个新的context，以便后续使用
	ctx = opentracing.ContextWithSpan(ctx, span)
	return handler(ctx, req)
}

func AccessLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	reqLog := "access log request: method %s, begin_time %d, request %v"
	beginTime := time.Now().Local().Unix()
	slog.Info(reqLog, info.FullMethod, beginTime, req)

	resp, err := handler(ctx, req)

	respLog := "access log response: method %s, begin_time %d, end_time %d, response %v"
	endTime := time.Now().Local().Unix()
	slog.Info(respLog, info.FullMethod, beginTime, endTime, resp)
	return resp, err
}

func ErrorLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		errLog := "error log: method %s, code %v, message %v, details %v"
		s := errcode.FromError(err)
		slog.Info(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
	}
	return resp, err
}

func Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer func() {
		if err := recover(); err != nil {
			recoveryLog := "recovery log: method %s, message %v, stack %s"
			slog.Info(recoveryLog, info.FullMethod, err, string(debug.Stack()[:]))
		}
	}()
	return handler(ctx, req)
}
