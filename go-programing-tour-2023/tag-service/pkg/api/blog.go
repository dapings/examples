package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/tag-service/global"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// TODO:
// 1. AccessToken 的缓存和刷新
// 2. 多 HTTP Method 的接口调用

const (
	AppKey    = "tour-2023"
	AppSecret = "tag-service"
)

type AccessToken struct {
	Token string `json:"token"`
}

type API struct {
	URL    string
	Client *http.Client
}

func NewAPI(url string) *API {
	return &API{URL: url, Client: http.DefaultClient}
}

func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	body, err := a.httpGet(ctx, fmt.Sprintf("%s?token=%s&name=%s", "api/v1/tags", token, name))
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 统一的 HTTP GET 请求方法
func (a *API) httpGet(ctx context.Context, urlPath string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", a.URL, urlPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// tracing
	// 解析上下文信息，创建并设置当前跨度的信息和标签内容(需传入上下文信息，保证链路完整性)
	// 传入附带的信息，并把它设置到对应的链路信息上；最后进行调用，并返回新的上下文，以便后续使用
	span, newCtx := opentracing.StartSpanFromContext(ctx,
		"HTTP GET: "+a.URL,
		opentracing.Tag{Key: string(ext.Component), Value: global.HTTPSpanTagVal},
	)
	span.SetTag("url", url)
	_ = opentracing.GlobalTracer().Inject(
		span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header),
	)

	req = req.WithContext(newCtx)
	client := http.Client{Timeout: 60 * time.Second} // default 60 * time.Second timeout
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			_ = body.Close()
		}
	}(resp.Body)
	defer span.Finish()

	// 读取消息实体，在实际封装中可以将其抽离
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// AccessToken 获取方法
func (a *API) getAccessToken(ctx context.Context) (string, error) {
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?app_key=%s&app_secret=%s", "auth", AppKey, AppSecret))
	if err != nil {
		return "", err
	}

	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return "", err
	}

	return accessToken.Token, nil
}
