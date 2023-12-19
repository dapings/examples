package sqldb

import (
	"go.uber.org/zap"
)

// 用 option 的模式注入参数。

type options struct {
	logger *zap.Logger
	sqlURL string
}

var defaultOptions = options{logger: zap.NewNop()}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithConnURL(sqlURL string) Option {
	return func(opts *options) {
		opts.sqlURL = sqlURL
	}
}
