package parse

import (
	"bytes"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
)

type Parser interface {
	ReadDocument(body []byte)
}

// 包括接收器的每个参数在输入函数/方法时被复制，返回时，对副本所做的更改将丢失。

type BaseParse struct {
	matchSyn string
}

func (p *BaseParse) HeaderSyn(headerReg string) {
	p.matchSyn = headerReg
}

func (p *BaseParse) ReadDocument(_ []byte) {}

type RegexpParse struct {
	BaseParse
}

func (p RegexpParse) ReadDocument(body []byte) {
	if len(body) == 0 || len(p.matchSyn) == 0 {
		return
	}

	// regexp.MustCompile函数会在编译时，提前解析好正则表达式内容，在一定程度上加速程序的运行。
	// [\s\S]*?，[\s\S] 任意字符串，*将前面任意字符匹配0次或无数次，?非贪婪匹配，找到第一次出现的地方，就认定匹配成功。
	// 由于回溯的原因，复杂的正则表达式，可能比较消耗CPU资源。
	// var headerReg = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)
	var headerReg = regexp.MustCompile(p.matchSyn)

	// 一个三维字节数组，第三层是字符实际对应的字节数组
	matches := headerReg.FindAllSubmatch(body, -1)
	for _, m := range matches {
		log.Println("fetch card news:", string(m[1][1]))
	}
}

type XPathParse struct {
	BaseParse
}

func (p XPathParse) ReadDocument(body []byte) {
	if len(body) == 0 || len(p.matchSyn) == 0 {
		return
	}

	// 解析HTML文本
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		log.Printf("htmlquery.Parse filed: %v", err)

		return
	}
	// 通过XPath语法查找符合条件的节点
	// `//div[@class="news_li"]/h2/a[@target="_blank"]`
	nodes := htmlquery.Find(doc, p.matchSyn)
	for _, node := range nodes {
		log.Println("fetch card news:", node.FirstChild.Data)
	}
}

type CSSSelection struct {
	BaseParse
}

func (p CSSSelection) ReadDocument(body []byte) {
	// 加载HTML文本
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Printf("load content filed: %v", err)

		return
	}
	// 根据CSS标签选择器的语法查找匹配的标签，并遍历输出a标签中的文本
	// "div.news_li h2 a[target=_blank]"
	doc.Find(p.matchSyn).Each(func(i int, s *goquery.Selection) {
		// 获取匹配的元素文本
		title := s.Text()
		log.Printf("review %d: %s\n", i, title)
	})
}
