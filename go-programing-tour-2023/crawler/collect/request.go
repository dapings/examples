package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type (
	Request struct {
		unique   string
		Task     *Task
		Url      string
		Method   string
		Depth    int
		Priority int
		RuleName string
	}

	ParseResult struct {
		Requests []*Request
		Items    []any
	}

	Task struct {
		Name        string // 用户界面显示的名称，且需保证唯一性
		Url         string
		Cookie      string
		WaitTime    time.Duration
		Reload      bool // 网站是否可以重复爬取
		MaxDepth    int
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
