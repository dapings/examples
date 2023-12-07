package main

import (
	"context"
	"log"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	// 创建新的metadata信息
	// md := metadata.New(map[string]string{"go": "programing", "tag": "service"})
	// newCtx := metadata.NewOutgoingContext(ctx, md)
	// 新增metadata信息
	newCtx := metadata.AppendToOutgoingContext(ctx, "tag-service", "go programing")
	clientConn, err := rpc.GetClientConn(newCtx, "127.0.0.1:8005", []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(
			middleware.UnaryCtxTimeout(),
			middleware.ClientTracing(),
		)),
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
