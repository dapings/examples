package spider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"regexp"
	"time"
)

type (
	Request struct {
		Task     *Task
		URL      string
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
	block := md5.Sum([]byte(r.URL + r.Method))

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
			URL:      u,
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
	if ok := re.Match(c.Body); !ok {
		return ParseResult{Items: make([]any, 0)}
	}

	return ParseResult{Items: []any{c.Req.URL}}
}

func (c *Context) Output(data any) *DataCell {
	res := &DataCell{Task: c.Req.Task}
	res.Data = make(map[string]any)
	res.Data["Task"] = c.Req.Task.Name
	res.Data["Rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["URL"] = c.Req.URL
	res.Data["Time"] = time.Now().Format(time.DateTime)

	return res
}

func (c *Context) GetRule(ruleName string) *Rule {
	return c.Req.Task.Rule.Trunk[ruleName]
}
