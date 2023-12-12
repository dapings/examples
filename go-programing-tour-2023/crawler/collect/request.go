package collect

import (
	"errors"
	"sync"
	"time"
)

type (
	Request struct {
		Task      *Task
		Url       string
		Depth     int
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
