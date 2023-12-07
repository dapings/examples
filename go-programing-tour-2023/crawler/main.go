package main

import (
	"bytes"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
)

var cssReg = "div.news_li h2 a[target=_blank]"

func main() {
	url := "https://www.thepaper.cn/"
	// url := "https://www.chinanews.com.cn/"
	// douban timeout
	// url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = collect.BrowserFetch{Timeout: 300 * time.Millisecond}
	body, err := f.Get(url)
	if err != nil {
		log.Printf("read content failed: %v", err)
		return
	}

	// 加载HTML文本
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Printf("load content filed: %v", err)
		return
	}
	// 根据CSS标签选择器的语法查找匹配的标签，并遍历输出a标签中的文本
	doc.Find(cssReg).Each(func(i int, s *goquery.Selection) {
		// 获取匹配的元素文本
		title := s.Text()
		log.Printf("review %d: %s\n", i, title)
	})

	collect.HandleLinks(body)
}
