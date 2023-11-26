package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/server"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 同端口上兼容多种协议：使用第三方开源库cmux来实现对多协议的支持。

var port string

func main() {
	l, err := runTCPServer(port)
	if err != nil {
		log.Fatalf("run tcp server err: %v", err)
	}

	mux := cmux.New(l)
	httpListener := mux.Match(cmux.HTTP1Fast())
	// cmux基于content-type头信息标识进行分流，grpc特定标识：application/grpc。
	grpcListener := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))

	grpcServer := runGRPCServer()
	httpServer := runHTTPServer(port)

	go grpcServer.Serve(grpcListener)
	go httpServer.Serve(httpListener)

	err = mux.Serve()
	if err != nil {
		log.Fatalf("run serve err: %v", err)
	}
}

func runTCPServer(port string) (net.Listener, error) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func runHTTPServer(port string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func runGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func init() {
	flag.StringVar(&port, "port", "8003", "启动端口")
	flag.Parse()
}
