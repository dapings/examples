package server

import (
	"context"
	"encoding/json"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	blogapi "github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/api"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/errcode"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TagServer struct {
	pb.UnimplementedTagServiceServer
	auth Auth
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	// 对服务进行内部调用，模拟类似gRPC服务内调的效果
	_, _ = t.internalGetTagList(ctx, r)

	if err := t.auth.check(ctx); err != nil {
		return nil, err
	}
	api := blogapi.NewAPI(global.BlogAddr)
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.ToGRPCError(errcode.ErrorGetTagListFail)
	}

	// JSON, Protobuf 结构体互转
	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)

	return &tagList, nil
}

type Auth struct {
	AppKey    string
	AppSecret string
}

// GetRequestMetadata 获取当前请求认证所需的元数据。
func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

// RequireTransportSecurity 是否需要基于TLS认证，进行安全传输。
func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func (a *Auth) getAppKey() string {
	return global.AppKey
}

func (a *Auth) getAppSecret() string {
	return global.AppSecret
}

func (a *Auth) check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)

	var appKey, appSecret string
	if val, ok := md["app_key"]; ok {
		appKey = val[0]
	}
	if val, ok := md["app_secret"]; ok {
		appSecret = val[0]
	}
	if appKey != a.getAppKey() || appSecret != a.getAppSecret() {
		return errcode.ToGRPCError(errcode.Unauthorized)
	}

	return nil
}

// 模拟类似gRPC服务内调的效果
func (t *TagServer) internalGetTagList(ctx context.Context, _ *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	clientConn, err := rpc.GetClientConn(ctx, global.TagServerAddr,
		// 客户端拦截器的相关注册行为是在调用grpc.Dial或grpc.DialContext之前，通过DialOption配置选项进行注册的。
		[]grpc.DialOption{grpc.WithUnaryInterceptor(
			grpcmiddleware.ChainUnaryClient(
				middleware.UnaryCtxTimeout(),
				middleware.ClientTracing(),
			),
		)})
	if err != nil {
		return nil, errcode.ToGRPCError(errcode.Fail)
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
		return nil, errcode.ToGRPCError(errcode.Fail)
	}

	return resp, nil
}
