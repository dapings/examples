package sqlstorage

import (
	"go.uber.org/zap"
)

type options struct {
	logger     *zap.Logger
	sqlUrl     string
	BatchCount int
}

var defaultOptions = options{logger: zap.NewNop()}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithSQLUrl(sqlUrl string) Option {
	return func(opts *options) {
		opts.sqlUrl = sqlUrl
	}
}

func WithBatchCount(batchCount int) Option {
	return func(opts *options) {
		opts.BatchCount = batchCount
	}
}