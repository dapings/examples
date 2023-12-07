package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var cssReg = "div.news_li h2 a[target=_blank]"

func main() {
	url := "https://www.thepaper.cn/"
	// url := "https://www.chinanews.com.cn/"
	body, err := Fetch(url)
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

	handleLinks(body)
}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("fetch url(%s) error: %v", url, err)
		panic(err)
	}

	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			_ = closer.Close()
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("error status code: %v", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	// 将HTML文本从特定编码转换为utf8编码
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

// DetermineEncoding 检测并返回当前HTML文本的编码格式
func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	// 如果返回的文本小于1024个字节，则说明文本有问题，直接使用utf8编码
	peek, err := r.Peek(1024)
	if err != nil {
		log.Printf("determin encoding err: %v", err)
		return unicode.UTF8
	}

	// 检测并返回当前HTML文本的编码格式
	e, _, _ := charset.DetermineEncoding(peek, "")
	return e
}

func handleLinks(body []byte) {
	if len(body) == 0 {
		return
	}

	numNews := bytes.Count(body, []byte("news_"))
	log.Printf("homepage has %d news class!\n", numNews)

	numLinks := strings.Count(string(body), "<a")
	log.Printf("homepage has %d links!\n", numLinks)

	numLinks = bytes.Count(body, []byte("<a"))
	log.Printf("homepage has %d links!\n", numLinks)

	exist := strings.Contains(string(body), "疫情")
	log.Printf("是否存在疫情:%v\n", exist)

	exist = bytes.Contains(body, []byte("疫情"))
	log.Printf("是否存在疫情:%v\n", exist)

	// log.Printf("body:%s", string(body))
}
