package cmd

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd/internal"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/master"
	"github.com/spf13/cobra"
	"go-micro.dev/v4/config"
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

	HTTPListenAddr string
	GRPCListenAddr string
	masterID       string
)

func init() {
	masterCmd.Flags().StringVar(&masterID, "id", "1", "set master id")
	masterCmd.Flags().StringVar(&HTTPListenAddr, "http_addr", ":8081", "set HTTP listen addr")
	masterCmd.Flags().StringVar(&GRPCListenAddr, "grpc_addr", ":9091", "set gRPC listen addr")
}

func RunMaster() {
	var (
		cfg     config.Config
		logger  *zap.Logger
		sconfig *internal.ServerConfig
		err     error
	)

	if cfg, err = internal.LoadConfig(); err != nil {
		panic(err)
	}

	if logger, err = internal.ConfigLogger(cfg); err != nil {
		panic(err)
	}

	if sconfig, err = internal.ConfigMasterServer(cfg, logger); err != nil {
		panic(err)
	}

	if _, err = master.New(
		masterID,
		master.WithLogger(logger.Named("master")),
		master.WithGRPCAddr(GRPCListenAddr),
		master.WithRegistryURL(sconfig.RegistryAddr),
	); err != nil {
		panic(err)
	}

	sconfig.ID = masterID
	sconfig.GRPCListenAddr = GRPCListenAddr
	sconfig.HTTPListenAddr = HTTPListenAddr

	// start http proxy to gRPC
	go internal.RunHTTPServer(logger, *sconfig)

	// start grpc server
	internal.RunGRPCServer(logger, *sconfig)
}
