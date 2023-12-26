package engine

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	WorkCount   int
	Fetcher     spider.Fetcher
	Storage     spider.Storage
	Logger      *zap.Logger
	Seeds       []*spider.Task
	registryURL string
	scheduler   Scheduler
}

var defaultOptions = options{Logger: zap.NewNop()}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.Logger = logger
	}
}

func WithFetcher(fetcher spider.Fetcher) Option {
	return func(opts *options) {
		opts.Fetcher = fetcher
	}
}

func WithWorkCount(workCount int) Option {
	return func(opts *options) {
		opts.WorkCount = workCount
	}
}

func WithSeeds(seeds []*spider.Task) Option {
	return func(opts *options) {
		opts.Seeds = seeds
	}
}

func WithScheduler(scheduler Scheduler) Option {
	return func(opts *options) {
		opts.scheduler = scheduler
	}
}

func WithStorage(s spider.Storage) Option {
	return func(opts *options) {
		opts.Storage = s
	}
}

func WithRegistryURL(registryURL string) Option {
	return func(opts *options) {
		opts.registryURL = registryURL
	}
}
