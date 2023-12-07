package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func main() {
	// url := "https://www.thepaper.cn/"
	url := "https://www.chinanews.com.cn/"
	body, err := Fetch(url)
	if err != nil {
		log.Printf("read content failed: %v", err)
		return
	}

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
