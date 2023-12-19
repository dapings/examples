package spider

import (
	"sync"
)

type (
	Property struct {
		Name     string `json:"name"` // 用户界面显示的名称，且需保证唯一性
		URL      string `json:"url"`
		Cookie   string `json:"cookie"`
		WaitTime int64  `json:"wait_time"` // 随机休眠时间，秒
		Reload   bool   `json:"reload"`    // 网站是否可以重复爬取
		MaxDepth int64  `json:"max_depth"`
	}

	LimitConfig struct {
		EventCount int
		EventDur   int // 秒
		Bucket     int // 桶大小
	}

	TaskConfig struct {
		Name     string
		Cookie   string
		WaitTime int64
		Reload   bool
		MaxDepth int64
		Fetcher  string
		Limits   []LimitConfig
	}

	// Task 一个任务实例
	Task struct {
		Visited     map[string]bool
		VisitedLock sync.Mutex
		Rule        RuleTree
		Options
	}
)

func NewTask(opts ...Option) *Task {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	d := &Task{}
	d.Options = options

	return d
}

type Fetcher interface {
	Get(req *Request) ([]byte, error)
}
