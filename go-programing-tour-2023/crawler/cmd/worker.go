package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd/internal"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/generator"
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

	WorkerServiceName = "go.micro.server.worker"

	WorkerHTTPListenAddr  string
	WorkerGRPCListenAddr  string
	WorkerPProfListenAddr string
	cluster               bool
	workerPodIP           string
	workerID              string
)

func init() {
	workerCmd.Flags().StringVar(&cfgFile, "config", "config.toml", "set config file")
	workerCmd.Flags().StringVar(&workerID, "id", "", "set worker id")
	workerCmd.Flags().StringVar(&workerPodIP, "podIP", "", "set worker pod ip")
	workerCmd.Flags().StringVar(&WorkerHTTPListenAddr, "http_addr", ":8080", "set HTTP listen addr")
	workerCmd.Flags().StringVar(&WorkerGRPCListenAddr, "grpc_addr", ":9090", "set gRPC listen addr")
	workerCmd.Flags().StringVar(&WorkerPProfListenAddr, "pprof_addr", ":9981", "set pprof listen addr")
	workerCmd.Flags().BoolVar(&cluster, "cluster", true, "run mode")
}

func RunWorker() {
	// start pprof.
	go func() {
		if err := http.ListenAndServe(WorkerPProfListenAddr, nil); err != nil {
			panic(err)
		}
	}()

	var (
		cfg       config.Config
		logger    *zap.Logger
		proxyFunc proxy.Func
		timeout   int
		fetcher   spider.Fetcher
		storager  spider.Storage
		sconfig   *internal.ServerConfig
		seeds     []*spider.Task
		s         *engine.Crawler
		err       error
	)

	if cfg, err = internal.LoadConfig(cfgFile); err != nil {
		panic(err)
	}

	if logger, err = internal.ConfigLogger(cfg); err != nil {
		panic(err)
	}

	logger.Named("worker")

	if proxyFunc, timeout, err = internal.ConfigProxyFunc(cfg, logger); err != nil {
		panic(err)
	}

	fetcher = internal.ConfigFetcher(proxyFunc, timeout, logger)

	if storager, err = internal.ConfigStorager(cfg, logger); err != nil {
		panic(err)
	}

	// init tasks
	if seeds, err = internal.ConfigTasks(cfg, fetcher, storager, logger); err != nil {
		panic(err)
	}

	if sconfig, err = internal.ConfigWorkerServer(cfg, logger); err != nil {
		panic(err)
	}

	if _, err = internal.ConfigWorkerEngine(sconfig, seeds, fetcher, storager, logger); err != nil {
		panic(err)
	}

	if workerID == "" {
		workerID = fmt.Sprintf("%d", time.Now().Local().UnixNano())

		if workerPodIP != "" {
			workerID = strconv.Itoa(int(generator.IDByIP(workerPodIP)))
		}
	}

	sconfig.ID = workerID
	sconfig.HTTPListenAddr = WorkerHTTPListenAddr
	sconfig.GRPCListenAddr = WorkerGRPCListenAddr
	WorkerServiceName = sconfig.Name

	// worker start
	id := sconfig.Name + "-" + workerID
	logger.Debug("worker id", zap.String("id", id))
	go s.Run(id, cluster)

	// start http proxy to gRPC
	go internal.RunHTTPServer(logger, *sconfig)

	// start grpc server
	internal.RunGRPCServer(nil, logger, *sconfig, nil)
}
