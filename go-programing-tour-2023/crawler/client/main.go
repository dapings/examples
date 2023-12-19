package main

import (
	"context"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	pb "github.com/dapings/examples/go-programing-tour-2023/crawler/protos/greeter"
	grpccli "github.com/go-micro/plugins/v4/client/grpc"
	etcdReg "github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// logger
	plugin := log.NewStdoutPlugin(zapcore.DebugLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init")

	// set zap global logger
	zap.ReplaceGlobals(logger)

	// option模式注入注册中心etcd的地址。
	reg := etcdReg.NewRegistry(registry.Addrs(":2379"))
	// 用option的模式注入参数；在默认情况下生成的服务器并不是gRPC类型的。
	service := micro.NewService(
		micro.Client(grpccli.NewClient()),
		// go-micro 注入etcd中的Key为/micro/registry/go.micro.server.worker/go.micro.server.worker-1
		micro.Registry(reg), // 注入register模块，用于指定注册中心，并定时发送自己的健康状况用于保活
		micro.WrapHandler(log.MicroServerWrapper(logger)),
	)
	service.Init()

	// 通过服务器的注册名
	cl := pb.NewGreeterService("go.micro.server.worker", service.Client())

	rsp, err := cl.Hello(context.Background(), &pb.Request{Name: "micro"})
	if err != nil {
		logger.Fatal("greeter hello failed", zap.Error(err))
	}
	logger.Info("greeter hello resp", zap.String("greeting", rsp.GetGreeting()))
}
