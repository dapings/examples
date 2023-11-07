package middleware

import (
	"context"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

// 在对所有路由方法调用之前生效，因此在注册路由行为之前进行注册。

func Tracing() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var newCtx context.Context
		var span opentracing.Span
		spanCtx, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				ctx.Request.Context(),
				global.Tracer,
				ctx.Request.URL.Path)
		} else {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				ctx.Request.Context(),
				global.Tracer,
				ctx.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"})
		}
		defer span.Finish()

		var traceID string
		var spanID string
		var spanContext = span.Context()
		switch spanContext.(type) {
		case jaeger.SpanContext:
			jaegerCtx := spanContext.(jaeger.SpanContext)
			traceID = jaegerCtx.TraceID().String()
			spanID = jaegerCtx.SpanID().String()
		}
		ctx.Set("X-Trace-ID", traceID)
		ctx.Set("X-Span-ID", spanID)
		ctx.Request = ctx.Request.WithContext(newCtx)

		ctx.Next()
	}
}
