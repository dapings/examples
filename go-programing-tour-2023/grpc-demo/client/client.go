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
