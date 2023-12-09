package main

import (
	"log"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
)

func main() {
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	proxyFunc, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		log.Printf("round robin proxy switcher failed: %v\n", err)
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
		log.Printf("read content failed: %v", err)
		return
	}

	collect.HandleLinks(body)

	var p parse.Parser = parse.CSSSelection{}
	var cssReg = "div.news_li h2 a[target=_blank]"

	p.WithHeaderSyn(cssReg)
	p.ReadDocument(body)
}
