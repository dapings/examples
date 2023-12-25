package cmd

import (
	"net/http"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd/internal"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/master"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/spf13/cobra"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

var (
	masterCmd = &cobra.Command{
		Use:   "master",
		Short: "run master service.",
		Long:  "run master service.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			RunMaster()
		},
	}

	MasterHTTPListenAddr  string
	MasterGRPCListenAddr  string
	MasterPProfListenAddr string
	masterID              string
)

func init() {
	masterCmd.Flags().StringVar(&masterID, "id", "1", "set master id")
	masterCmd.Flags().StringVar(&MasterHTTPListenAddr, "http_addr", ":8081", "set HTTP listen addr")
	masterCmd.Flags().StringVar(&MasterGRPCListenAddr, "grpc_addr", ":9091", "set gRPC listen addr")
	masterCmd.Flags().StringVar(&MasterPProfListenAddr, "pprof_addr", ":9081", "set pprof listen addr")
}

func RunMaster() {
	// start pprof.
	go func() {
		if err := http.ListenAndServe(MasterPProfListenAddr, nil); err != nil {
			panic(err)
		}
	}()

	var (
		cfg     config.Config
		logger  *zap.Logger
		seeds   []*spider.Task
		sconfig *internal.ServerConfig
		m       *master.Master
		err     error
	)

	if cfg, err = internal.LoadConfig(); err != nil {
		panic(err)
	}

	if logger, err = internal.ConfigLogger(cfg); err != nil {
		panic(err)
	}

	if seeds, err = internal.ConfigTasks(cfg, nil, nil, logger); err != nil {
		panic(err)
	}

	if sconfig, err = internal.ConfigMasterServer(cfg, logger); err != nil {
		panic(err)
	}

	reg := etcd.NewRegistry(registry.Addrs(sconfig.RegistryAddr))

	if m, err = master.New(
		masterID,
		master.WithLogger(logger.Named("master")),
		master.WithGRPCAddr(MasterGRPCListenAddr),
		master.WithRegistryURL(sconfig.RegistryAddr),
		master.WithRegistry(reg),
		master.WithSeeds(seeds),
	); err != nil {
		logger.Error("init master failed", zap.Error(err))

		panic(err)
	}

	sconfig.ID = masterID
	sconfig.GRPCListenAddr = MasterGRPCListenAddr
	sconfig.HTTPListenAddr = MasterHTTPListenAddr

	// start http proxy to gRPC
	go internal.RunHTTPServer(logger, *sconfig)

	// start grpc server
	internal.RunGRPCServer(m, logger, *sconfig, reg)
}
