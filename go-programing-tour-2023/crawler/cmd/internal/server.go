package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/master"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/protos/crawler"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/protos/greeter"
	grpccli "github.com/go-micro/plugins/v4/client/grpc"
	etcdReg "github.com/go-micro/plugins/v4/registry/etcd"
	gs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, resp *greeter.Response) error {
	resp.Greeting = "Hello " + req.Name

	return nil
}

type ServerConfig struct {
	GRPCListenAddr   string
	HTTPListenAddr   string
	ID               string
	RegistryAddr     string
	RegisterTTL      int
	RegisterInterval int
	Name             string
	ClientTimeOut    int
}

func RunGRPCServer(hdlr *master.Master, logger *zap.Logger, cfg ServerConfig, reg registry.Registry) {
	// start grpc server
	if reg == nil {
		// option模式注入注册中心etcd的地址。
		reg = etcdReg.NewRegistry(registry.Addrs(cfg.RegistryAddr))
	}
	// 用option的模式注入参数；在默认情况下生成的服务器并不是gRPC类型的。
	service := micro.NewService(
		// gRPC插件生成一个gRPC Server
		micro.Server(gs.NewServer(server.Id(cfg.ID))), // 指定特殊的ID来替换随机的ID
		micro.Address(cfg.GRPCListenAddr),
		micro.Name(cfg.Name), // 服务器的名字
		// go-micro 注入etcd中的Key为/micro/registry/go.micro.server.worker/go.micro.server.worker-1
		micro.Registry(reg), // 注入register模块，用于指定注册中心，并定时发送自己的健康状况用于保活
		micro.WrapHandler(log.MicroServerWrapper(logger)),
		micro.RegisterTTL(time.Duration(cfg.RegisterTTL)*time.Second),
		micro.RegisterInterval(time.Duration(cfg.RegisterInterval)*time.Second),
		micro.Client(grpccli.NewClient()),
	)

	// 设置micro客户端默认超时时间为10秒
	if err := service.Client().Init(client.RequestTimeout(time.Duration(cfg.ClientTimeOut) * time.Second)); err != nil {
		logger.Sugar().Error("micro client inti error. ", zap.String("error:", err.Error()))

		return
	}

	service.Init()

	// register handler
	var err error
	if hdlr != nil {
		hdlr.SetForwardCli(crawler.NewCrawlerMasterService(cfg.Name, service.Client()))
		err = crawler.RegisterCrawlerMasterHandler(service.Server(), hdlr)
	} else {
		// default greeter.
		err = greeter.RegisterGreeterHandler(service.Server(), new(Greeter))
	}

	if err != nil {
		logger.Fatal("register handler failed", zap.Error(err))
	}

	if err := service.Run(); err != nil {
		logger.Fatal("grpc server stop", zap.Error(err))
	}
}

func RunHTTPServer(logger *zap.Logger, cfg ServerConfig) {
	if logger == nil {
		logger = zap.L()
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 指定要转发到那个gRPC服务器。
	err := crawler.RegisterCrawlerMasterGwFromEndpoint(ctx, mux, cfg.GRPCListenAddr, opts)
	if err != nil {
		logger.Fatal("pb register gw backend from ep(grpc server endpoint) failed", zap.Error(err))
	}

	err = greeter.RegisterGreeterGwFromEndpoint(ctx, mux, cfg.GRPCListenAddr, opts)
	if err != nil {
		logger.Fatal("pb register gw backend from ep(grpc server endpoint) failed", zap.Error(err))
	}

	zap.S().Debugf("start http server listening on %v proxy to grpc server;%v", cfg.HTTPListenAddr, cfg.GRPCListenAddr)

	if err := http.ListenAndServe(cfg.HTTPListenAddr, mux); err != nil {
		logger.Fatal("http listen and serve failed", zap.Error(err))
	}
}
