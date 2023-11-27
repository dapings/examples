package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/naming"
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

	r := &naming.GRPCResolver{Client: cli}
	target := fmt.Sprintf("/etcdv3://go-programming-tour/grpc/%s", global.ServiceName)
	opts := []grpc.DialOption{grpc.WithBalancer(grpc.RoundRobin(r)), grpc.WithBlock()}

	ctx := context.Background()
	clientConn, err := rpc.GetClientConn(ctx, target, opts)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer clientConn.Close()
	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}
