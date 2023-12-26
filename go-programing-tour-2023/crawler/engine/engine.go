package engine

import (
	"runtime/debug"
	"sync"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/master"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Crawler struct {
	id          string
	out         chan spider.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex

	failures    map[string]*spider.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex

	resources map[string]*master.ResourceSpec
	rlock     sync.Mutex

	etcdCli *clientv3.Client
	options
}

func NewEngine(opts ...Option) (*Crawler, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	e := &Crawler{}

	e.Visited = make(map[string]bool, 100)
	e.out = make(chan spider.ParseResult)
	e.failures = make(map[string]*spider.Request)
	e.options = options

	// 任务添加默认的采集器与存储器
	for _, task := range Store.list {
		task.Fetcher = e.Fetcher
		task.Storage = e.Storage
	}

	endpoints := []string{e.registryURL}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}
	e.etcdCli = cli

	return e, nil
}

func (e *Crawler) Run(id string, cluster bool) {
	e.id = id
	if !cluster {
		e.handleSeeds()
	}

	go func() {
		err := e.loadResource()
		if err != nil {
			_ = e.loadResource()
		}
	}()
	go e.watchResource()
	go e.Schedule()

	for i := 0; i < e.WorkCount; i++ {
		go e.CreateWork()
	}

	e.HandleResult()
}

func (e *Crawler) Schedule() {
	go e.scheduler.Schedule()
}

func (e *Crawler) CreateWork() {
	defer func() {
		if err := recover(); err != nil {
			e.Logger.Error("worker panic", zap.Any("err", err), zap.String("stack", string(debug.Stack())))
		}
	}()

	for {
		req := e.scheduler.Pull()
		if err := req.Check(); err != nil {
			e.Logger.Error("check failed", zap.Error(err))

			continue
		}

		if !req.Task.Reload && e.HasVisited(req) {
			e.Logger.Debug("request has visited", zap.String("url:", req.URL))

			continue
		}

		e.StoreVisited(req)

		body, err := req.Fetch()
		if err != nil {
			e.Logger.Error("read content failed",
				zap.Int("len", len(body)), zap.Error(err), zap.String("url", req.URL))

			e.SetFailure(req)

			continue
		}

		if len(body) < 6000 {
			e.Logger.Error("read content failed",
				zap.Int("len", len(body)), zap.String("url", req.URL))

			e.SetFailure(req)

			continue
		}

		e.Logger.Info("get content", zap.Int("len", len(body)))

		rule := req.Task.Rule.Trunk[req.RuleName]

		result, parsedErr := rule.ParseFunc(&spider.Context{
			Body: body,
			Req:  req,
		})

		if parsedErr != nil {
			e.Logger.Error("rule.ParseFunc failed", zap.Error(parsedErr))

			continue
		}

		if len(result.Requests) > 0 {
			go e.scheduler.Push(result.Requests...)
		}

		e.out <- result
	}
}

func (e *Crawler) HandleResult() {
	for result := range e.out {
		for _, item := range result.Items {
			switch d := item.(type) {
			case *spider.DataCell:
				if err := d.Task.Storage.Save(d); err != nil {
					// TODO: when store error, skip or other handle method ?
					e.Logger.Error("task.Storage.Save failed", zap.Error(err))

					continue
				}
			}

			e.Logger.Sugar().Info("get result", item)
		}
	}
}

func (e *Crawler) HasVisited(r *spider.Request) bool {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()

	unique := r.Unique()

	return e.Visited[unique]
}

func (e *Crawler) StoreVisited(reqs ...*spider.Request) {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()

	for _, r := range reqs {
		unique := r.Unique()
		e.Visited[unique] = true
	}
}

func (e *Crawler) SetFailure(req *spider.Request) {
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
	// TODO: 失败2次，加载到失败队列中?
}
