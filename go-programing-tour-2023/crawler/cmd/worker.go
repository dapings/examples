package cmd

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd/internal"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"go-micro.dev/v4/config"
	"go.uber.org/zap"
)

func RunWorker() {
	var (
		cfg       config.Config
		logger    *zap.Logger
		proxyFunc proxy.Func
		timeout   int
		fetcher   spider.Fetcher
		storager  spider.Storage
		sconfig   *internal.ServerConfig
		s         *engine.Crawler
		err       error
	)

	if cfg, err = internal.LoadConfig(); err != nil {
		panic(err)
	}

	if logger, err = internal.ConfigLogger(cfg); err != nil {
		panic(err)
	}

	if proxyFunc, timeout, err = internal.ConfigProxyFunc(cfg, logger); err != nil {
		panic(err)
	}

	fetcher = internal.ConfigFetcher(proxyFunc, timeout, logger)

	if storager, err = internal.ConfigStorager(cfg, logger); err != nil {
		panic(err)
	}

	// init tasks
	if s, err = internal.ConfigTasks(cfg, fetcher, storager, logger); err != nil {
		panic(err)
	}

	if sconfig, err = internal.ConfigWorkerServer(cfg, logger); err != nil {
		panic(err)
	}

	// worker start
	go s.Run()

	// start http proxy to gRPC
	go internal.RunHTTPServer(logger, *sconfig)

	// start grpc server
	internal.RunGRPCServer(logger, *sconfig)
}
