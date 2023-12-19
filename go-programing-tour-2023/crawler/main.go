package main

import (
	"context"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/limiter"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubanbook"
	pb "github.com/dapings/examples/go-programing-tour-2023/crawler/protos/greeter"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/sqldb"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/storage"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/storage/sqlstorage"
	etcdReg "github.com/go-micro/plugins/v4/registry/etcd"
	gs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// logger
	plugin := log.NewStdoutPlugin(zapcore.DebugLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init")

	// set zap global logger
	zap.ReplaceGlobals(logger)

	var err error

	// proxy
	var proxyFunc proxy.Func

	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}

	if proxyFunc, err = proxy.RoundRobinProxySwitcher(proxyURLs...); err != nil {
		logger.Error("round robin proxy switcher failed", zap.Error(err))

		return
	}

	// url := "https://www.thepaper.cn/"
	// url := "https://www.chinanews.com.cn/"
	// url := "https://google.com.hk"

	// douban timeout
	// url := "https://book.douban.com/subject/1007305/"
	// fetcher
	var f collect.Fetcher = collect.BrowserFetch{Timeout: 300 * time.Millisecond, Proxy: proxyFunc, Logger: logger}

	// storage
	var storager storage.Storage

	if storager, err = sqlstorage.New(
		sqlstorage.WithSQLURL(sqldb.ConStrWithMySQL),
		sqlstorage.WithLogger(logger.Named("SQLDB")),
		sqlstorage.WithBatchCount(2),
	); err != nil {
		logger.Error("create storage failed", zap.Error(err))

		return
	}

	// speed limiter
	// 2秒1个
	secondLimit := rate.NewLimiter(limiter.Per(1, 2*time.Second), 1)
	// 60秒20个
	minuteLimit := rate.NewLimiter(limiter.Per(20, 1*time.Minute), 20)
	multiLimiter := limiter.Multi(secondLimit, minuteLimit)

	// init tasks
	// seeds slice cap
	var seeds = make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			Name: doubanbook.BookListTaskName,
		},
		Fetcher: f,
		Storage: storager,
		Limit:   multiLimiter,
	})

	s := engine.NewEngine(
		engine.WithWorkCount(5),
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	// worker start
	go s.Run()

	// start http proxy to gRPC
	go RunHTTPServer()

	RunGRPCServer(logger)
}

func RunGRPCServer(logger *zap.Logger) {
	// start grpc server
	// option模式注入注册中心etcd的地址。
	reg := etcdReg.NewRegistry(registry.Addrs(":2379"))
	// 用option的模式注入参数；在默认情况下生成的服务器并不是gRPC类型的。
	service := micro.NewService(
		// gRPC插件生成一个gRPC Server
		micro.Server(gs.NewServer(server.Id("1"))), // 指定特殊的ID来替换随机的ID
		micro.Address(":9090"),
		micro.Name("go.micro.server.worker"), // 服务器的名字
		// go-micro 注入etcd中的Key为/micro/registry/go.micro.server.worker/go.micro.server.worker-1
		micro.Registry(reg), // 注入register模块，用于指定注册中心，并定时发送自己的健康状况用于保活
		micro.WrapHandler(log.MicroServerWrapper(logger)),
		micro.RegisterTTL(60*time.Second),
		micro.RegisterInterval(15*time.Second),
	)

	// 设置micro客户端默认超时时间为10秒
	if err := service.Client().Init(client.RequestTimeout(10 * time.Second)); err != nil {
		logger.Sugar().Error("micro client inti error. ", zap.String("error:", err.Error()))

		return
	}

	service.Init()

	if err := pb.RegisterGreeterHandler(service.Server(), new(Greeter)); err != nil {
		logger.Fatal("register handler failed", zap.Error(err))
	}

	if err := service.Run(); err != nil {
		logger.Fatal("grpc server stop", zap.Error(err))
	}
}

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	resp.Greeting = "Hello " + req.Name

	return nil
}

func RunHTTPServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 指定要转发到那个gRPC服务器。
	err := pb.RegisterGreeterGwFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		zap.L().Fatal("pb register gw from ep failed", zap.Error(err))
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		zap.L().Fatal("http listen and serve failed", zap.Error(err))
	}
}
