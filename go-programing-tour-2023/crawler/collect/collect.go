package collect

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type BaseFetch struct{}

func (BaseFetch) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("fetch url(%s) error: %v", url, err)
		return nil, err
	}

	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			_ = closer.Close()
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("error status code: %v", resp.StatusCode)
		return nil, fmt.Errorf("error status code: %v", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	// 将HTML文本从特定编码转换为utf8编码
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
}

// Get 模拟浏览器访问
func (b BrowserFetch) Get(url string) ([]byte, error) {
	client := &http.Client{Timeout: b.Timeout}
	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("http new request filed: %v\n", err)
		return nil, fmt.Errorf("http new request filed: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("fetch url(%s) error: %v", url, err)
		return nil, err
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

func HandleLinks(body []byte) {
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
