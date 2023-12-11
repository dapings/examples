package collect

import (
	"errors"
	"time"
)

type (
	Request struct {
		Url       string
		Cookie    string
		WaitTime  time.Duration
		Depth     int
		MaxDepth  int
		ParseFunc func([]byte, *Request) ParseResult
	}

	ParseResult struct {
		Requests []*Request
		Items    []any
	}
)

func (r *Request) Check() error {
	if r.Depth > r.MaxDepth {
		return errors.New("max depth limit reached")
	}

	return nil
}
