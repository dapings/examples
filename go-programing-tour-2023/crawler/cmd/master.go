package cmd

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd/internal"
	"go-micro.dev/v4/config"
	"go.uber.org/zap"
)

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

	// start http proxy to gRPC
	go internal.RunHTTPServer(logger, *sconfig)

	// start grpc server
	internal.RunGRPCServer(logger, *sconfig)
}
