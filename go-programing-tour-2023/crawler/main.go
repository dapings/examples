package main

import (
	"log"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse"
)

func main() {
	url := "https://www.thepaper.cn/"
	// url := "https://www.chinanews.com.cn/"
	// url := "https://google.com"

	// douban timeout
	// url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = collect.BrowserFetch{Timeout: 300 * time.Millisecond}
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
