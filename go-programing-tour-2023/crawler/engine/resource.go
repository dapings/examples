package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/master"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (e *Crawler) handleSeeds() {
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

	go e.scheduler.Push(reqQueue...)
}

func (e *Crawler) watchResource() {
	watch := e.etcdCli.Watch(context.Background(), master.ResourcePath, clientv3.WithPrefix())
	for w := range watch {
		if w.Err() != nil {
			e.Logger.Error("watch resource failed", zap.Error(w.Err()))

			continue
		}

		if w.Canceled {
			e.Logger.Error("watch resource canceled")

			return
		}

		for _, ev := range w.Events {
			spec, err := master.Decode(ev.Kv.Value)
			if err != nil {
				e.Logger.Error("decode etcd value failed", zap.Error(err))

				continue
			}

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					e.Logger.Info("receive create resource", zap.Any("spec", spec))
				}

				if ev.IsModify() {
					e.Logger.Info("receive update resource", zap.Any("spec", spec))
				}

				e.runTasks(spec.Name)
			case clientv3.EventTypeDelete:
				e.Logger.Info("receive delete resource", zap.Any("spec", spec))
			}
		}
	}
}

func (e *Crawler) loadResource() error {
	resp, err := e.etcdCli.Get(context.Background(), master.ResourcePath, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return fmt.Errorf("get resource from etcd failed: %w", err)
	}

	rs := make(map[string]*master.ResourceSpec)
	for _, kv := range resp.Kvs {
		r, err := master.Decode(kv.Value)
		if err == nil && r != nil {
			id := getID(r.AssignedNode)
			if len(id) > 0 && e.id == id {
				rs[r.Name] = r
			}
		}
	}

	e.Logger.Info("leader init load resource", zap.Int("length", len(rs)))

	e.rlock.Lock()
	defer e.rlock.Unlock()

	e.resources = rs

	for _, r := range e.resources {
		e.runTasks(r.Name)
	}

	return nil
}

func (e *Crawler) runTasks(taskName string) {
	t, ok := Store.Hash[taskName]
	if !ok {
		e.Logger.Error("not found preset tasks", zap.String("task name", taskName))

		return
	}

	reqs, err := t.Rule.Root()
	if err != nil {
		e.Logger.Error("get task rule tree root failed", zap.String("task name", taskName), zap.Error(err))

		return
	}

	for _, req := range reqs {
		req.Task = t
	}

	e.scheduler.Push(reqs...)
}

func getID(assignedNode string) string {
	s := strings.Split(assignedNode, "|")
	if len(s) < 2 {
		return ""
	}

	return s[0]
}
