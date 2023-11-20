package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/dapings/examples/go-programing-tour-2023/grpc-demo/protos"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "the server port")
)

func main() {
	flag.Parsed()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", l.Addr())

	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("received: %v", in.GetName())
	return &pb.HelloReply{
		Message: "hello " + in.GetName(),
	}, nil
}
