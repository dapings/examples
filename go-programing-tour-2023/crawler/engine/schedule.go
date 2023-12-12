package engine

import (
	"sync"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"go.uber.org/zap"
)

type Crawler struct {
	out         chan collect.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	options
}

type Scheduler interface {
	Schedule()
	Push(...*collect.Request)
	Pull() *collect.Request
}

type Schedule struct {
	reqChan    chan *collect.Request
	workerChan chan *collect.Request
	reqQueue   []*collect.Request
	Logger     *zap.Logger
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	e := &Crawler{}
	e.Visited = make(map[string]bool, 100)
	e.out = make(chan collect.ParseResult)
	e.options = options
	return e
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	s.reqChan = make(chan *collect.Request)
	s.workerChan = make(chan *collect.Request)

	return s
}

func (e *Crawler) Run() {
	go e.Schedule()

	for i := 0; i < e.WorkCount; i++ {
		go e.CreateWork()
	}

	e.HandleResult()
}

func (s *Schedule) Push(reqQueue ...*collect.Request) {
	for _, req := range reqQueue {
		s.reqChan <- req
	}
}

func (s *Schedule) Pull() *collect.Request {
	return <-s.workerChan
}

func (s *Schedule) Output() *collect.Request {
	return <-s.workerChan
}

func (s *Schedule) Schedule() {
	for {
		var req *collect.Request
		var ch chan *collect.Request

		if len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerChan
		}
		select {
		case r := <-s.reqChan:
			s.reqQueue = append(s.reqQueue, r)
		case ch <- req:
		}
	}
}

func (e *Crawler) Schedule() {
	var reqQueue []*collect.Request
	for _, seed := range e.Seeds {
		seed.RootRequest.Task = seed
		seed.RootRequest.Url = seed.Url
		reqQueue = append(reqQueue, seed.RootRequest)
	}
	go e.scheduler.Schedule()
	go e.scheduler.Push(reqQueue...)
}

func (e *Crawler) CreateWork() {
	for {
		r := e.scheduler.Pull()
		if err := r.Check(); err != nil {
			e.Logger.Error("check failed", zap.Error(err))
			continue
		}
		if e.HasVisited(r) {
			e.Logger.Debug("request has visited", zap.String("url:", r.Url))
			continue
		}
		e.StoreVisited(r)

		body, err := r.Task.Fetcher.Get(r)
		if len(body) < 6000 {
			e.Logger.Error("read content failed",
				zap.Int("len", len(body)), zap.String("url", r.Url))
			continue
		}
		if err != nil {
			e.Logger.Error("read content failed",
				zap.Error(err), zap.String("url", r.Url))
			continue
		}
		e.Logger.Info("get content", zap.Int("len", len(body)))
		result := r.ParseFunc(body, r)

		if len(result.Requests) > 0 {
			go e.scheduler.Push(result.Requests...)
		}

		e.out <- result
	}
}

func (e *Crawler) HandleResult() {
	for {
		select {
		case result := <-e.out:
			for _, item := range result.Items {
				// TODO: store
				e.Logger.Sugar().Info("get result", item)
			}
		}
	}
}

func (e *Crawler) HasVisited(r *collect.Request) bool {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()
	unique := r.Unique()
	return e.Visited[unique]
}

func (e *Crawler) StoreVisited(reqs ...*collect.Request) {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()

	for _, r := range reqs {
		unique := r.Unique()
		e.Visited[unique] = true
	}
}
