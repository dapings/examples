package main

import (
	"fmt"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubangroup"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
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
	var seeds []*collect.Request
	for i := 0; i <= 0; i += 25 {
		seeds = append(seeds, &collect.Request{
			Url:       fmt.Sprintf(doubangroup.DiscussionURL, i),
			Cookie:    doubangroup.Cookie,
			ParseFunc: doubangroup.ParseURL,
		})
	}

	s := engine.ScheduleEngine{
		WorkCount: 5,
		Fetcher:   f,
		Logger:    logger,
		Seeds:     seeds,
	}
	s.Run()
}
