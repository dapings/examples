package collect

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"regexp"
	"sync"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/limiter"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/storage"
	"go.uber.org/zap"
)

type (
	Request struct {
		unique   string
		Task     *Task
		Url      string
		Method   string
		Depth    int64
		Priority int64
		RuleName string
		TmpData  *Temp
	}

	ParseResult struct {
		Requests []*Request
		Items    []any
	}

	Property struct {
		Name     string `json:"name"` // 用户界面显示的名称，且需保证唯一性
		Url      string `json:"url"`
		Cookie   string `json:"cookie"`
		WaitTime int64  `json:"wait_time"`
		Reload   bool   `json:"reload"` // 网站是否可以重复爬取
		MaxDepth int64  `json:"max_depth"`
	}

	Task struct {
		Property
		Visited     map[string]bool
		VisitedLock sync.Mutex
		Fetcher     Fetcher
		Storage     storage.Storage
		Rule        RuleTree
		Logger      *zap.Logger
		Limit       limiter.RateLimiter
	}

	Context struct {
		Body []byte
		Req  *Request
	}
)

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("max depth limit reached")
	}

	return nil
}

func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}

func (r *Request) Fetch() ([]byte, error) {
	if err := r.Task.Limit.Wait(context.Background()); err != nil {
		return nil, err
	}
	// 随机休眠，模拟人类行为
	sleepTime := rand.Int63n(r.Task.WaitTime * 1000)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	return r.Task.Fetcher.Get(r)
}

func (c *Context) ParseJSReg(name, reg string) ParseResult {
	re := regexp.MustCompile(reg)
	matches := re.FindAllSubmatch(c.Body, -1)
	result := ParseResult{}
	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &Request{
			Task:     c.Req.Task,
			Url:      u,
			Method:   "GET",
			Depth:    c.Req.Depth + 1,
			Priority: 0,
			RuleName: name,
		})
	}
	return result
}

func (c *Context) OutputJS(reg string) ParseResult {
	re := regexp.MustCompile(reg)
	ok := re.Match(c.Body)
	if !ok {
		return ParseResult{Items: make([]any, 0)}
	}
	return ParseResult{Items: []any{c.Req.Url}}
}

func (c *Context) Output(data any) *storage.DataCell {
	res := &storage.DataCell{}
	res.Data = make(map[string]any)
	res.Data["Task"] = c.Req.Task.Name
	res.Data["Rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["Url"] = c.Req.Url
	res.Data["Time"] = time.Now().Format(time.DateTime)
	return res
}

func (c *Context) GetRule(ruleName string) *Rule {
	return c.Req.Task.Rule.Trunk[ruleName]
}
