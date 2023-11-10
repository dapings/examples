package main

import (
	"context"
	"log"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/tracer"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	newCtx := metadata.AppendToOutgoingContext(ctx, "", "")
	clientConn, err := rpc.GetClientConn(newCtx, "127.0.0.1:8004",
		[]grpc.DialOption{grpc.WithUnaryInterceptor(
			grpcmiddleware.ChainUnaryClient(
				middleware.UnaryCtxTimeout(),
				middleware.ClientTracing(),
			),
		)})
	if err != nil {
		log.Fatalf("rpc.GetClientConn err: %v", err)
	}
	defer func(clientConn *grpc.ClientConn) {
		err := clientConn.Close()
		if err != nil {
			_ = clientConn.Close()
		}
	}(clientConn)

	// 业务逻辑：查询标签列表
	// TODO
}

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setTracer err: %v", err)
	}
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("tag-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
