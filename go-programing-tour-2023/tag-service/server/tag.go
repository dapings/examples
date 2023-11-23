package server

import (
	"context"
	"encoding/json"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	blogapi "github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/api"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/errcode"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
)

type TagServer struct {
	pb.UnimplementedTagServiceServer
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
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
