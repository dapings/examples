package engine

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubanbook"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubangroup"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/parse/doubangroupjs"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/robertkrimen/otto"
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
