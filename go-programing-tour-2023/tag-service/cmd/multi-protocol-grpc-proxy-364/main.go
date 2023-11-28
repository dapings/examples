package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/proxy/grpcproxy"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var port string

func main() {
	err := runServer(port)
	if err != nil {
		log.Fatalf("run serve err: %v", err)
	}
}

func runServer(port string) error {
	httpMux := runHTTPServer()
	gRPCServer := runGRPCServer()

	endpoint := "0.0.0.0:" + port
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{rpc.GetGRPCDialOptionWithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	httpMux.Handle("/", gwmux)

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: time.Second * 60,
	})
	if err != nil {
		return err
	}
	defer etcdClient.Close()

	target := fmt.Sprintf("/etcdv3://go-programming-tour/grpc/%s", global.ServiceName)
	grpcproxy.Register(etcdClient, target, ":"+port, 60)

	return http.ListenAndServe(":"+port, gRPCHandlerFunc(gRPCServer, httpMux))
}

func runHTTPServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	return serveMux
}

func runGRPCServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.AccessLog,
			middleware.ErrorLog,
			middleware.Recovery,
		)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func gRPCHandlerFunc(gRPCServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			gRPCServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func gRPCGatewayERROR(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler,
	w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	httpError := httpError{
		Code:    int32(s.Code()),
		Message: s.Message(),
	}
	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*pb.Error); ok {
			httpError.Code = v.Code
			httpError.Message = v.Message
		}
	}

	resp, _ := json.Marshal(httpError)
	w.Header().Set("Content-Type", marshaler.ContentType(""))
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	_, _ = w.Write(resp)
}

func init() {
	flag.StringVar(&port, "port", "8004", "启动端口")
	flag.Parse()
}

type httpError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
