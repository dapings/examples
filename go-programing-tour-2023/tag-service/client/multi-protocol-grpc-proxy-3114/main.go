package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/naming"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	"google.golang.org/grpc"
)

func main() {
	config := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: time.Second * 60,
	}
	cli, err := clientv3.New(config)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	// creates a grpc.Watcher for a target to track its resolution changes.
	// etcd SDK对官方Resolver,Watcher两个接口的具体实现。
	r := &naming.GRPCResolver{Client: cli}
	target := fmt.Sprintf("/etcdv3://go-programming-tour/grpc/%s", global.ServiceName)
	// 设置负载均衡器(策略：round robin)和连接建立的要求(Block: 达到Ready状态才返回)。
	opts := []grpc.DialOption{grpc.WithBalancer(grpc.RoundRobin(r)), grpc.WithBlock()}

	ctx := context.Background()
	clientConn, err := rpc.GetClientConn(ctx, target, opts)
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
	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}
