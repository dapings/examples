package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 不同的两个端口，监听不同协议的流量：使用两个协程分别监听 http endpoint, grpc endpoint 实际是一个阻塞行为。

var (
	grpcPort string
	httpPort string
)

func main() {
	errChan := make(chan error)
	go func() {
		err := runHttpServer(httpPort)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		err := runGRPCServer(grpcPort)
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Fatalf("run server err: %v", err)
	}
}

func runHttpServer(port string) error {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	return http.ListenAndServe(":"+port, serverMux)
}

func runGRPCServer(port string) error {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	err = s.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "8001", "gRPC启动端口")
	flag.StringVar(&httpPort, "http_port", "8002", "HTTP启动端口")
	flag.Parse()
}
