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
		unique    string
		Task      *Task
		Url       string
		Method    string
		Depth     int
		Priority  int
		ParseFunc func([]byte, *Request) ParseResult
	}

	ParseResult struct {
		Requests []*Request
		Items    []any
	}

	Task struct {
		Url         string
		Cookie      string
		WaitTime    time.Duration
		MaxDepth    int
		Visited     map[string]bool
		VisitedLock sync.Mutex
		RootRequest *Request
		Fetcher     Fetcher
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
