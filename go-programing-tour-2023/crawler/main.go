package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func main() {
	// logger
	plugin := log.NewStdoutPlugin(zapcore.DebugLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init")

	// set zap global logger
	zap.ReplaceGlobals(logger)

	// proxy
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	var err error
	var proxyFunc proxy.ProxyFunc
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
		sqlstorage.WithSQLUrl(sqldb.ConStrWithMySQL),
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
	multiLimiter := limiter.MultiLimiter(secondLimit, minuteLimit)

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
	go HandleHTTP()

	// start grpc server
	reg := etcdReg.NewRegistry(registry.Addrs(":2379"))
	service := micro.NewService(
		micro.Server(gs.NewServer(server.Id("1"))),
		micro.Address(":9090"),
		micro.Registry(reg),
		micro.Name("go.micro.server.worker"),
	)
	service.Init()
	_ = pb.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err := service.Run(); err != nil {
		logger.Fatal("grpc server stop")
	}
}

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	resp.Greeting = "Hello " + req.Name
	return nil
}

func HandleHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterGreeterGwFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		fmt.Println(err)
	}
	_ = http.ListenAndServe(":8080", mux)
}
