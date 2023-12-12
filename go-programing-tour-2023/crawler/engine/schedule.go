package engine

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"go.uber.org/zap"
)

type Schedule struct {
	reqChan    chan *collect.Request
	workerChan chan *collect.Request
	out        chan collect.ParseResult
	options
}

func NewSchedule(opts ...Option) *Schedule {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &Schedule{}
	s.options = options
	return s
}

func (s *Schedule) Run() {
	s.reqChan = make(chan *collect.Request)
	s.workerChan = make(chan *collect.Request)
	s.out = make(chan collect.ParseResult)

	go s.Schedule()

	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}

	s.HandleResult()
}

func (s *Schedule) Schedule() {
	var reqQueue []*collect.Request
	for _, seed := range s.Seeds {
		seed.RootRequest.Task = seed
		seed.RootRequest.Url = seed.Url
		reqQueue = append(reqQueue, seed.RootRequest)
	}
	go func() {
		for {
			var req *collect.Request
			var ch chan *collect.Request

			if len(reqQueue) > 0 {
				req = reqQueue[0]
				reqQueue = reqQueue[1:]
				ch = s.workerChan
			}
			select {
			case r := <-s.reqChan:
				reqQueue = append(reqQueue, r)
			case ch <- req:
			}
		}
	}()
}

func (s *Schedule) CreateWork() {
	for {
		r := <-s.workerChan
		if err := r.Check(); err != nil {
			s.Logger.Error("check failed", zap.Error(err))
			continue
		}
		body, err := s.Fetcher.Get(r)
		if len(body) < 6000 {
			s.Logger.Error("read content failed",
				zap.Int("len", len(body)), zap.String("url", r.Url))
			continue
		}
		if err != nil {
			s.Logger.Error("read content failed",
				zap.Error(err), zap.String("url", r.Url))
			continue
		}
		s.Logger.Info("get content", zap.Int("len", len(body)))
		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

func (s *Schedule) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, request := range result.Requests {
				s.reqChan <- request
			}

			for _, item := range result.Items {
				// TODO: store
				s.Logger.Sugar().Info("get result", item)
			}
		}
	}
}
