package cmd

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd/internal"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/spf13/cobra"
	"go-micro.dev/v4/config"
	"go.uber.org/zap"
)

var (
	workerCmd = &cobra.Command{
		Use:   "worker",
		Short: "run worker service.",
		Long:  "run worker service.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			RunWorker()
		},
	}

	WorkerHTTPListenAddr string
	WorkerGRPCListenAddr string
	workerID             string
)

func init() {
	workerCmd.Flags().StringVar(&workerID, "id", "1", "set worker id")
	workerCmd.Flags().StringVar(&WorkerHTTPListenAddr, "http_addr", ":8080", "set HTTP listen addr")
	workerCmd.Flags().StringVar(&WorkerGRPCListenAddr, "grpc_addr", ":9090", "set gRPC listen addr")
}

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

	sconfig.ID = workerID
	sconfig.HTTPListenAddr = WorkerHTTPListenAddr
	sconfig.GRPCListenAddr = WorkerGRPCListenAddr

	// worker start
	// go s.Run()

	// start http proxy to gRPC
	go internal.RunHTTPServer(logger, *sconfig)

	// start grpc server
	internal.RunGRPCServer(logger, *sconfig)
}
