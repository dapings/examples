package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/limiter"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	pb "github.com/dapings/examples/go-programing-tour-2023/crawler/protos/greeter"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/sqldb"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/storage/sqlstorage"
	"github.com/go-micro/plugins/v4/config/encoder/toml"
	etcdReg "github.com/go-micro/plugins/v4/registry/etcd"
	gs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source"
	"go-micro.dev/v4/config/source/file"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var err error
	// load config
	enc := toml.NewEncoder()
	cfg, cfgErr := config.NewConfig(config.WithReader(json.NewReader(reader.WithEncoder(enc))))
	if cfgErr != nil {
		panic(cfgErr)
	}
	err = cfg.Load(file.NewSource(file.WithPath("config.toml"), source.WithEncoder(enc)))
	if err != nil {
		panic(err)
	}

	logText := cfg.Get("logLevel").String("INFO")
	logLevel, err := zapcore.ParseLevel(logText)
	if err != nil {
		return
	}

	// logger
	plugin := log.NewStdoutPlugin(logLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init")

	// set zap global logger
	zap.ReplaceGlobals(logger)

	// proxy
	var proxyFunc proxy.Func

	proxyURLs := cfg.Get("fetcher", "proxy").StringSlice([]string{})
	timeout := cfg.Get("fetcher", "timeout").Int(5000)
	logger.Sugar().Info("proxy list:", proxyURLs, " timeout: ", timeout)

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
	var f spider.Fetcher = collect.BrowserFetch{Timeout: time.Duration(timeout) * time.Millisecond, Proxy: proxyFunc, Logger: logger}

	// storage
	var storager spider.Storage

	sqlURL := cfg.Get("storage", "sqlURL").String("")
	sqldb.ConStrWithMySQL = sqlURL

	if storager, err = sqlstorage.New(
		sqlstorage.WithSQLURL(sqldb.ConStrWithMySQL),
		sqlstorage.WithLogger(logger.Named("sqlDB")),
		sqlstorage.WithBatchCount(2),
	); err != nil {
		logger.Error("create sql storage failed", zap.Error(err))

		return
	}

	// init tasks
	// seeds slice cap
	var seeds = make([]*spider.Task, 0, 1000)
	var taskCfg []spider.TaskConfig
	if err := cfg.Get("Tasks").Scan(&taskCfg); err != nil {
		logger.Fatal("init seed tasks failed", zap.Error(err))
	}
	seeds = ParseTaskConfig(logger, f, storager, taskCfg)

	s := engine.NewEngine(
		engine.WithWorkCount(5),
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	// worker start
	go s.Run()

	var sconfig ServerConfig
	if err := cfg.Get("GRPCServer").Scan(&sconfig); err != nil {
		logger.Fatal("get gRPC Server config failed", zap.Error(err))
	}
	logger.Sugar().Debugf("grpc server config, %+v", sconfig)

	// start http proxy to gRPC
	go RunHTTPServer(sconfig)

	RunGRPCServer(logger, sconfig)
}

func ParseTaskConfig(logger *zap.Logger, f spider.Fetcher, s spider.Storage, cfgs []spider.TaskConfig) []*spider.Task {
	tasks := make([]*spider.Task, 0, 1000)
	for _, cfg := range cfgs {
		t := spider.NewTask(
			spider.WithName(cfg.Name),
			spider.WithReload(cfg.Reload),
			spider.WithCookie(cfg.Cookie),
			spider.WithLogger(logger),
			spider.WithStorage(s),
		)

		if cfg.WaitTime > 0 {
			t.WaitTime = cfg.WaitTime
		}

		if cfg.MaxDepth > 0 {
			t.MaxDepth = cfg.MaxDepth
		}

		var limits []limiter.RateLimiter
		if len(cfg.Limits) > 0 {
			for _, limitCfg := range cfg.Limits {
				// speed limiter 限速器，2秒1个，60秒20个
				l := rate.NewLimiter(limiter.Per(limitCfg.EventCount, time.Duration(limitCfg.EventDur)*time.Second), 1)
				limits = append(limits, l)
			}
			multiLimiter := limiter.Multi(limits...)
			t.Limit = multiLimiter
		}

		switch cfg.Fetcher {
		case "browser":
			t.Fetcher = f
		}

		tasks = append(tasks, t)
	}

	return tasks
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

func RunGRPCServer(logger *zap.Logger, cfg ServerConfig) {
	// start grpc server
	// option模式注入注册中心etcd的地址。
	reg := etcdReg.NewRegistry(registry.Addrs(cfg.RegistryAddr))
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
	)

	// 设置micro客户端默认超时时间为10秒
	if err := service.Client().Init(client.RequestTimeout(time.Duration(cfg.ClientTimeOut) * time.Second)); err != nil {
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

func RunHTTPServer(cfg ServerConfig) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 指定要转发到那个gRPC服务器。
	err := pb.RegisterGreeterGwFromEndpoint(ctx, mux, cfg.GRPCListenAddr, opts)
	if err != nil {
		zap.L().Fatal("pb register gw from ep failed", zap.Error(err))
	}

	zap.S().Debugf("start http server listening on %v proxy to grpc server;%v", cfg.HTTPListenAddr, cfg.GRPCListenAddr)

	if err := http.ListenAndServe(cfg.HTTPListenAddr, mux); err != nil {
		zap.L().Fatal("http listen and serve failed", zap.Error(err))
	}
}
