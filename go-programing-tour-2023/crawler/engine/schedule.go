package engine

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"go.uber.org/zap"
)

type ScheduleEngine struct {
	reqChan    chan *collect.Request
	workerChan chan *collect.Request
	WorkCount  int
	Fetcher    collect.Fetcher
	Logger     *zap.Logger
	out        chan collect.ParseResult
	Seeds      []*collect.Request
}

func (s *ScheduleEngine) Run() {
	s.reqChan = make(chan *collect.Request)
	s.workerChan = make(chan *collect.Request)
	s.out = make(chan collect.ParseResult)

	go s.Schedule()

	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}

	s.HandleResult()
}

func (s *ScheduleEngine) Schedule() {
	var reqQueue = s.Seeds
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

func (s *ScheduleEngine) CreateWork() {
	for {
		r := <-s.workerChan
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("read content failed", zap.Error(err))
			continue
		}
		s.Logger.Info("get content", zap.Int("len", len(body)))
		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

func (s *ScheduleEngine) HandleResult() {
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
