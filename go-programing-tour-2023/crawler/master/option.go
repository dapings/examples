package master

import (
	"go.uber.org/zap"
)

type options struct {
	logger      *zap.Logger
	registryURL string
	GRPCAddr    string
}

type Option func(opts *options)

var defaultOptions = options{logger: zap.NewNop()}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithRegistryURL(registryURL string) Option {
	return func(opts *options) {
		opts.registryURL = registryURL
	}
}

func WithGRPCAddr(addr string) Option {
	return func(opts *options) {
		opts.GRPCAddr = addr
	}
}
