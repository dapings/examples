package main

import (
	"context"
	"log"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/tracer"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	newCtx := metadata.AppendToOutgoingContext(ctx, "tag-server", "Go Programing")
	clientConn, err := rpc.GetClientConn(newCtx, global.TagServerAddr,
		// 客户端拦截器的相关注册行为是在调用grpc.Dial或grpc.DialContext之前，通过DialOption配置选项进行注册的。
		[]grpc.DialOption{grpc.WithUnaryInterceptor(
			grpcmiddleware.ChainUnaryClient(
				middleware.UnaryCtxTimeout(),
				middleware.ClientTracing(),
				// 重试功能：确定是否需要重试的维度是以错误码为标准的，
				// 若是，需要明确状态码的规则，确保多服务的状态码的标准是一致的(通过基础框架、公共库等方式落地)，
				// 另外，要尽可能保证接口设计是幂等的，保证重试不会造成灾难性的问题，如重复扣库存等。
				grpc_retry.UnaryClientInterceptor(
					grpc_retry.WithMax(5),
					grpc_retry.WithCodes(
						codes.Unknown,
						codes.Internal,
						codes.DeadlineExceeded,
					),
				),
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
	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
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
