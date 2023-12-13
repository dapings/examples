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

	"github.com/dapings/examples/go-programing-tour-2023/crawler/extensions"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(*Request) ([]byte, error)
}

type BaseFetch struct{}

func (BaseFetch) Get(request *Request) ([]byte, error) {
	resp, err := http.Get(request.Url)
	if err != nil {
		log.Printf("fetch url(%s) error: %v", request.Url, err)
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
	Logger  *zap.Logger
}

// Get 模拟浏览器访问
func (b BrowserFetch) Get(request *Request) ([]byte, error) {
	client := &http.Client{Timeout: b.Timeout}
	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}
	req, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		b.Logger.Error("http new request filed: ", zap.Error(err))
		return nil, fmt.Errorf("http new request filed: %v", err)
	}
	if len(request.Task.Cookie) > 0 {
		req.Header.Set("Cookie", request.Task.Cookie)
	}

	req.Header.Set("User-Agent", extensions.GenerateRandomUA())

	resp, err := client.Do(req)
	time.Sleep(request.Task.WaitTime)
	if err != nil {
		b.Logger.Error("fetch url error: ", zap.String("fetch url", request.Url), zap.Error(err))
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
