package main

import (
	"context"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/rpc"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/swagger"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/tracer"
	pb "github.com/dapings/examples/go-programing-tour-2023/tag-service/protos"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/server"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var port string

func main() {
	err := runServer(port)
	if err != nil {
		log.Fatalf("run server err: %v", err)
	}
}

func runServer(port string) error {
	httpMux := runHTTPServer()
	grpcServer := runGRPCServer()
	gatewayMux := runGRPCGatewayServer()

	httpMux.Handle("/", gatewayMux)
	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcServer, httpMux))
}

// 创建一个新的 http 多路复用器
func runHTTPServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	// 访问本地生成的 .swagger.json，查看所生成的API描述信息
	serveMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}
		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("proto", p)

		http.ServeFile(w, r, p)
	})

	return serveMux
}

func runGRPCGatewayServer() *gwruntime.ServeMux {
	ctx := context.Background()
	endpoint := "0.0.0.0:" + port
	var serverMuxOpts []gwruntime.ServeMuxOption
	// 使用 grpc-gateway 中的 runtime.WithMarshalerOption方法，注册所需要的MIME类型及对应的Marshaler，
	// 默认使用 application/json，application/jsonpb 进行序列化。
	serverMuxOpts = append(serverMuxOpts,
		gwruntime.WithMarshalerOption("*", &gwruntime.HTTPBodyMarshaler{Marshaler: &gwruntime.JSONPb{
			MarshalOptions:   protojson.MarshalOptions{EmitUnpopulated: true},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}}),
	)
	gwmux := gwruntime.NewServeMux(serverMuxOpts...)
	// gRPC Server/Client 在启动和调用时，必须明确其是否加密，`DialOpton grpc.WithInsecure`指定为非加密模式。
	dopts := []grpc.DialOption{rpc.GetGRPCDialOptionWithInsecure()}
	// 注册方法事件，内部会自动转换并拨号到grpc endpoint，并在上下文结束后关闭连接。
	// 主要进行gRPC连接的创建和管控。
	// 1. 将当前RPC方法预定义的HTTP Endpoint注册到传入的gwmux HTTP多路复用器中
	// 2. 超时时间通过 ctx 进行控制
	// 3. 请求/响应数据，根据传入的 MIME 类型进行序列化，默认是 json
	// 4. 将gRPC metadata 转换为 context，便于使用
	_ = pb.RegisterTagServiceHandlerFromEndpoint(ctx, gwmux, endpoint, dopts)

	return gwmux
}

func runGRPCServer() *grpc.Server {
	// grpc.ServerOption 设置Server的相关属性，如credentials,keepalive等参数。
	// 拦截器在此注册，需以指定的类型封装，如一元拦截器的类型必须为grpc.UnaryInterceptor。
	// 在grpc.UnaryInterceptor中嵌套grpc_middleware.ChainUnaryServer，以链式方式达到用多个拦截器的目的。
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			middleware.AccessLog,
			middleware.ErrorLog,
			middleware.Recovery,
			middleware.ServerTracing,
		)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	// 注册gRPC反射服务，才要以使用 grpcurl 工具调试gRPC接口
	reflection.Register(s)
	return s
}

// grpc(HTTP/2)和HTTP/1.1通过Header中的Content-Type和ProtoMajor标识进行分流
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	// h2c标识允许通过明文TCP运行HTTP/2协议，用于HTTP/1.1升级标头字段和标识HTTP/2 over TCP
	// 官方标准库golang.org/x/net/http2/h2c 实现了HTTP/2的未加密模式，可直接使用
	// h2c.NewHandler内部逻辑是拦截所有h2c流量，根据不同的请求流量类型(Content-Type)，将其劫持并重定向到相应的Handler中去处理，
	// 最终完成在同个端口上既能提供HTTP/1.1的功能，又能提供HTTP/2的功能。
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 对gRPC(HTTP/2),HTTP/1.1的流量区分
		// ProtoMajor 代表客户端请求的版本号，客户端始终使用HTTP/1.1或HTTP/2协议
		// 通过gRPC的标志位application/grpc，确定流量的类型
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setTracer err: %v", err)
	}
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("tag-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
