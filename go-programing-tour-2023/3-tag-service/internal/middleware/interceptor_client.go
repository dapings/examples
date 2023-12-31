package middleware

import (
	"context"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/meta"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientTracing() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var (
			parentSpanCtx opentracing.SpanContext
			spanOpts      []opentracing.StartSpanOption
			// 解析上下文信息
			parentSpan = opentracing.SpanFromContext(ctx)
		)
		// 检查 是否包含上一级的跨度信息
		// 若存在，则获取上一级的上下文信息，作为本次跨度的父级
		if parentSpan != nil {
			parentSpanCtx = parentSpan.Context()
			spanOpts = append(spanOpts, opentracing.ChildOf(parentSpanCtx))
		}
		spanOpts = append(spanOpts, []opentracing.StartSpanOption{
			// 创建和设置本次跨度的标签信息
			opentracing.Tag{Key: string(ext.Component), Value: global.GRPCSpanTagVal},
			ext.SpanKindRPCClient,
		}...)

		span := global.Tracer.StartSpan(method, spanOpts...)
		defer span.Finish()

		// 对传出的md信息进行转换，把它设置到新的上下文信息中，以便后续使用
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		_ = global.Tracer.Inject(span.Context(), opentracing.TextMap, meta.MetadataTextMap{MD: md})
		newCtx := opentracing.ContextWithSpan(metadata.NewOutgoingContext(ctx, md), span)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

// 超时控制(对上下文超时时间的设置和适当控制)，是在微服务架构中非常重要的一个保命项。

func UnaryCtxTimeout() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := defaultCxtTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamCtxTimeout() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx, cancel := defaultCxtTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func defaultCxtTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	// 检查，若未设置截止时间，则返回false；设置默认超时时间
	if _, ok := ctx.Deadline(); !ok {
		defaultTimeout := 60 * time.Second
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx, cancel
}
