package engine

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"go.uber.org/zap"
)

type Config struct {
	WorkCount int
	Fetcher   collect.Fetcher
	Logger    *zap.Logger
	Seeds     []*collect.Request
}
