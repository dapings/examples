package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"regexp"
	"sync"
	"time"
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
	}

	ParseResult struct {
		Requests []*Request
		Items    []any
	}

	Property struct {
		Name     string        `json:"name"` // 用户界面显示的名称，且需保证唯一性
		Url      string        `json:"url"`
		Cookie   string        `json:"cookie"`
		WaitTime time.Duration `json:"wait_time"`
		Reload   bool          `json:"reload"` // 网站是否可以重复爬取
		MaxDepth int64         `json:"max_depth"`
	}

	Task struct {
		Property
		Visited     map[string]bool
		VisitedLock sync.Mutex
		Fetcher     Fetcher
		Rule        RuleTree
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
