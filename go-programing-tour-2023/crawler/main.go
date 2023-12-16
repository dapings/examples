package main

import (
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/collector"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/collector/sqlstorage"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubanbook"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/sqldb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// logger
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init")

	// proxy
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	proxyFunc, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("round robin proxy switcher failed", zap.Error(err))
		return
	}

	// url := "https://www.thepaper.cn/"
	// url := "https://www.chinanews.com.cn/"
	// url := "https://google.com.hk"

	// douban timeout
	// url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = collect.BrowserFetch{Timeout: 300 * time.Millisecond, Proxy: proxyFunc, Logger: logger}
	// storage
	var storage collector.Storage
	storage, err = sqlstorage.New(
		sqlstorage.WithSQLUrl(sqldb.ConStrWithMySQL),
		sqlstorage.WithLogger(logger.Named("SQLDB")),
		sqlstorage.WithBatchCount(2),
	)
	if err != nil {
		logger.Error("create storage failed", zap.Error(err))
		return
	}
	// seeds slice cap
	var seeds = make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			Name: doubanbook.BookListTaskName,
		},
		Fetcher: f,
		Storage: storage,
	})

	s := engine.NewEngine(
		engine.WithWorkCount(5),
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)
	s.Run()
}
