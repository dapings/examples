package global

import (
	"github.com/opentracing/opentracing-go"
)

var (
	Tracer opentracing.Tracer

	HTTPSpanTagVal = "HTTP"
	GRPCSpanTagVal = "gRPC"
)
