package server

import (
	"context"
	"encoding/json"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	blogapi "github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/api"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/errcode"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
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
