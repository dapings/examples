package engine

import (
	"runtime/debug"
	"sync"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubanbook"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubangroup"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubangroupjs"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/robertkrimen/otto"
	"go.uber.org/zap"
)

func init() {
	Store.Add(doubanbook.DoubanBookTask)
	Store.Add(doubangroup.DoubangroupTask)
	Store.AddJSTask(doubangroupjs.DoubangroupJSTask)
}

// Store 全局爬虫(蜘蛛)任务实例
var Store = &CrawlerStore{
	list: make([]*spider.Task, 0),
	Hash: make(map[string]*spider.Task),
}

// GetFields 获取任务规则的配置项。
func GetFields(taskName, ruleName string) []string {
	return Store.Hash[taskName].Rule.Trunk[ruleName].ItemFields
}

type Crawler struct {
	out         chan spider.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	failures    map[string]*spider.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex
	options
}

type CrawlerStore struct {
	list []*spider.Task
	Hash map[string]*spider.Task
}

func (c *CrawlerStore) Add(task *spider.Task) {
	c.Hash[task.Name] = task
	c.list = append(c.list, task)
}

func (c *CrawlerStore) AddJSTask(m *spider.TaskModel) {
	task := &spider.Task{
		// Property: m.Property,
	}
	task.Rule.Root = func() ([]*spider.Request, error) {
		// allocate a new JavaScript runtime
		vm := otto.New()
		if err := vm.Set("AddJsReq", spider.AddJsReq); err != nil {
			return nil, err
		}

		v, evalErr := vm.Eval(m.Root)

		if evalErr != nil {
			return nil, evalErr
		}

		e, exportErr := v.Export()

		if exportErr != nil {
			return nil, exportErr
		}

		return e.([]*spider.Request), nil
	}

	for _, r := range m.Rules {
		parseFunc := func(parse string) func(ctx *spider.Context) (spider.ParseResult, error) {
			return func(ctx *spider.Context) (spider.ParseResult, error) {
				// allocate a new JavaScript runtime
				vm := otto.New()
				if err := vm.Set("ctx", ctx); err != nil {
					return spider.ParseResult{}, err
				}

				v, evalErr := vm.Eval(parse)

				if evalErr != nil {
					return spider.ParseResult{}, evalErr
				}

				e, exportErr := v.Export()

				if exportErr != nil {
					return spider.ParseResult{}, exportErr
				}

				if e == nil {
					return spider.ParseResult{}, exportErr
				}

				return e.(spider.ParseResult), exportErr
			}
		}(r.ParseFunc)

		if task.Rule.Trunk == nil {
			task.Rule.Trunk = make(map[string]*spider.Rule, 0)
		}

		task.Rule.Trunk[r.Name] = &spider.Rule{ParseFunc: parseFunc}
	}

	c.Hash[task.Name] = task
	c.list = append(c.list, task)
}

type Scheduler interface {
	Schedule()
	Push(reqQueue ...*spider.Request)
	Pull() *spider.Request
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	e := &Crawler{}

	e.Visited = make(map[string]bool, 100)
	e.out = make(chan spider.ParseResult)
	e.failures = make(map[string]*spider.Request)
	e.options = options

	return e
}

type Schedule struct {
	reqChan     chan *spider.Request
	workerChan  chan *spider.Request
	priReqQueue []*spider.Request
	reqQueue    []*spider.Request
	Logger      *zap.Logger
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	s.reqChan = make(chan *spider.Request)
	s.workerChan = make(chan *spider.Request)

	return s
}

func (s *Schedule) Push(reqQueue ...*spider.Request) {
	for _, req := range reqQueue {
		s.reqChan <- req
	}
}

func (s *Schedule) Pull() *spider.Request {
	return <-s.workerChan
}

func (s *Schedule) Schedule() {
	var req *spider.Request

	var ch chan *spider.Request

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
	var reqQueue []*spider.Request

	for _, task := range e.Seeds {
		t, ok := Store.Hash[task.Name]
		if !ok {
			e.Logger.Error("not find preset tasks", zap.String("task name", task.Name))
		}
		// task.Logger = e.Logger
		task.Rule = t.Rule
		rootReqs, err := task.Rule.Root()

		if err != nil {
			e.Logger.Error("get root failed", zap.Error(err))

			continue
		}

		for _, req := range rootReqs {
			req.Task = task
		}

		reqQueue = append(reqQueue, rootReqs...)
	}

	go e.scheduler.Schedule()
	go e.scheduler.Push(reqQueue...)
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
					// TODO: when store error, skip or other handle method
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
	// TODO: 失败2次，加载到失败队列中
}
