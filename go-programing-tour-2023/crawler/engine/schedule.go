package engine

import (
	"sync"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubangroup"
	"go.uber.org/zap"
)

func init() {
	Store.Add(doubangroup.DoubangroupTask)
}

var Store = &CrawlerStore{
	list: make([]*collect.Task, 0),
	hash: make(map[string]*collect.Task),
}

type Crawler struct {
	out         chan collect.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	failures    map[string]*collect.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex
	options
}

type CrawlerStore struct {
	list []*collect.Task
	hash map[string]*collect.Task
}

func (c *CrawlerStore) Add(task *collect.Task) {
	c.hash[task.Name] = task
	c.list = append(c.list, task)
}

type Scheduler interface {
	Schedule()
	Push(...*collect.Request)
	Pull() *collect.Request
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	e := &Crawler{}
	e.Visited = make(map[string]bool, 100)
	e.out = make(chan collect.ParseResult)
	e.failures = make(map[string]*collect.Request)
	e.options = options
	return e
}

type Schedule struct {
	reqChan     chan *collect.Request
	workerChan  chan *collect.Request
	priReqQueue []*collect.Request
	reqQueue    []*collect.Request
	Logger      *zap.Logger
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	s.reqChan = make(chan *collect.Request)
	s.workerChan = make(chan *collect.Request)

	return s
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
	var req *collect.Request
	var ch chan *collect.Request
	for {
		if req == nil && len(s.priReqQueue) > 0 {
			req = s.priReqQueue[0]
			s.priReqQueue = s.priReqQueue[1:]
			ch = s.workerChan
		}
		if req == nil && len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerChan
		}

		select {
		case r := <-s.reqChan:
			if r.Priority > 0 {
				s.priReqQueue = append(s.priReqQueue, r)
			} else {
				s.reqQueue = append(s.reqQueue, r)
			}
		case ch <- req:
			req = nil
			ch = nil
		}
	}
}

func (e *Crawler) Run() {
	go e.Schedule()

	for i := 0; i < e.WorkCount; i++ {
		go e.CreateWork()
	}

	e.HandleResult()
}

func (e *Crawler) Schedule() {
	var reqQueue []*collect.Request
	for _, seed := range e.Seeds {
		task := Store.hash[seed.Name]
		task.Fetcher = seed.Fetcher
		rootReqs := task.Rule.Root()
		for _, req := range rootReqs {
			req.Task = task
		}
		reqQueue = append(reqQueue, rootReqs...)
	}
	go e.scheduler.Schedule()
	go e.scheduler.Push(reqQueue...)
}

func (e *Crawler) CreateWork() {
	for {
		req := e.scheduler.Pull()
		if err := req.Check(); err != nil {
			e.Logger.Error("check failed", zap.Error(err))
			continue
		}
		if !req.Task.Reload && e.HasVisited(req) {
			e.Logger.Debug("request has visited", zap.String("url:", req.Url))
			continue
		}
		e.StoreVisited(req)

		body, err := req.Task.Fetcher.Get(req)
		if err != nil {
			e.Logger.Error("read content failed",
				zap.Int("len", len(body)), zap.Error(err), zap.String("url", req.Url))
			e.SetFailure(req)
			continue
		}
		if len(body) < 6000 {
			e.Logger.Error("read content failed",
				zap.Int("len", len(body)), zap.String("url", req.Url))
			e.SetFailure(req)
			continue
		}

		e.Logger.Info("get content", zap.Int("len", len(body)))
		rule := req.Task.Rule.Trunk[req.RuleName]
		result := rule.ParseFunc(&collect.Context{
			Body: body,
			Req:  req,
		})

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

func (e *Crawler) SetFailure(req *collect.Request) {
	if !req.Task.Reload {
		e.VisitedLock.Lock()
		unique := req.Unique()
		delete(e.Visited, unique)
		e.VisitedLock.Unlock()
	}

	e.failureLock.Unlock()
	defer e.failureLock.Unlock()
	if _, ok := e.failures[req.Unique()]; !ok {
		// 首次失败，再重新执行一次
		e.failures[req.Unique()] = req
		e.scheduler.Push(req)
	}
	// TODO: 失败2次，加载到失败队列中
}
