package main

import (
	"io"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	log2 "github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	plugin, c := log2.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer func(c io.Closer) {
		err := c.Close()
		if err != nil {
			_ = c.Close()
		}
	}(c)
	logger := log2.NewLogger(plugin)
	logger.Info("log init")

	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	proxyFunc, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("round robin proxy switcher failed", zap.Error(err))
		return
	}
	url := "https://www.thepaper.cn/"
	// url := "https://www.chinanews.com.cn/"
	// url := "https://google.com.hk"

	// douban timeout
	// url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = collect.BrowserFetch{Timeout: 300 * time.Millisecond, Proxy: proxyFunc}
	body, err := f.Get(url)
	if err != nil {
		logger.Error("read content failed", zap.Error(err))
		return
	}

	logger.Info("get content", zap.Int("len", len(body)))

	collect.HandleLinks(body)

	var p parse.Parser = parse.CSSSelection{}
	var cssReg = "div.news_li h2 a[target=_blank]"

	p.WithHeaderSyn(cssReg)
	p.ReadDocument(body)
}
