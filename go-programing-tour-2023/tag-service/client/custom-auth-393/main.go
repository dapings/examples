package main

import (
	"context"
	"log"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/server"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	auth := server.Auth{
		AppKey:    global.AppKey,
		AppSecret: global.AppSecret,
	}
	ctx := context.Background()
	// 创建新的metadata信息
	md := metadata.New(map[string]string{"go": "programing", "tag": "service"})
	newCtx := metadata.NewOutgoingContext(ctx, md)
	clientConn, err := rpc.GetClientConn(newCtx, "127.0.0.1:8005", []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(
			middleware.UnaryCtxTimeout(),
			middleware.ClientTracing(),
		)),
		// 注册自定义认证信息
		grpc.WithPerRPCCredentials(&auth),
	})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer func(clientConn *grpc.ClientConn) {
		err := clientConn.Close()
		if err != nil {
			_ = clientConn.Close()
		}
	}(clientConn)

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}
