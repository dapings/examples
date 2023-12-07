package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/dapings/examples/go-programing-tour-2023/grpc-demo/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "dapings world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the addr to connect to")
	name = flag.String("name", defaultName, "name to greet")
)

func main() {
	flag.Parse()

	// Dial 创建客户端连接，不是马上建立可用连接
	// 如果需要立刻打通与服务端的连接，需要设置WithBlock DialOption，这样当发起连接时会阻塞等待连接完成，使最终连接到达Ready状态
	// grpc.WithInsecure()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			_ = conn.Close()
		}
	}(conn)

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("not greet: %v", err)
	}
	log.Printf("greeting: %s", resp.GetMessage())
}
