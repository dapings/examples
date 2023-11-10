package main

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/tag-service/pkg/tracer"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/reflection"

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

func runHTTPServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		// Asset:     swagger.Asset,
		// AssetDir:  swagger.AssetDir,
		Prefix: "third_party/swagger-ui",
	})
	serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
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
	// ctx := context.Background()
	// endpoint := "0.0.0.0:"+port
	gwmux := gwruntime.NewServeMux()
	// dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	return gwmux
}

func runGRPCServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			middleware.AccessLog,
			middleware.ErrorLog,
			middleware.Recovery,
			middleware.ServerTracing,
		)),
	}
	s := grpc.NewServer(opts...)
	reflection.Register(s)
	return s
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
