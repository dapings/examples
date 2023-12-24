package master

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

type options struct {
	logger      *zap.Logger
	registryURL string
	GRPCAddr    string
	registry    registry.Registry
	Seeds       []*spider.Task
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

func WithRegistry(reg registry.Registry) Option {
	return func(opts *options) {
		opts.registry = reg
	}
}

func WithSeeds(seeds []*spider.Task) Option {
	return func(opts *options) {
		opts.Seeds = seeds
	}
}
