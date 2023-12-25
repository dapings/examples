package master

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const (
	ResourcePath = "/resources"

	MSG_ADD Command = iota
	MSG_DEL
)

type Command int

type Message struct {
	Cmd   Command
	Specs []*ResourceSpec
}

type ResourceSpec struct {
	ID           string
	Name         string
	AssignedNode string
	CreateTime   int64
}

func getResourcePath(name string) string {
	return fmt.Sprintf("%s/%s", ResourcePath, name)
}

func encode(spec *ResourceSpec) string {
	b, _ := json.Marshal(spec)
	return string(b)
}

func decode(ds []byte) (*ResourceSpec, error) {
	var s *ResourceSpec
	err := json.Unmarshal(ds, &s)
	return s, err
}

func (m *Master) Assign(r *ResourceSpec) (*registry.Node, error) {
	for _, n := range m.workerNodes {
		return n, nil
	}

	return nil, errors.New("no  worker nodes")
}

func (m *Master) AddSeed() {
	rs := make([]*ResourceSpec, 0, len(m.Seeds))
	for _, seed := range m.Seeds {
		resp, err := m.etcdCli.Get(context.Background(), getResourcePath(seed.Name), clientv3.WithSerializable())
		if err != nil {
			m.logger.Error("etcd get '"+seed.Name+"' resource failed", zap.Error(err))

			continue
		}

		if len(resp.Kvs) == 0 {
			r := &ResourceSpec{Name: seed.Name}
			rs = append(rs, r)
		}
	}

	m.AddResource(rs)
}

func (m *Master) AddResource(rs []*ResourceSpec) {
	for _, r := range rs {
		// r.ID =
		r.CreateTime = time.Now().Local().UnixNano()
	}
}

func (m *Master) loadResource() error {
	resp, err := m.etcdCli.Get(context.Background(), ResourcePath, clientv3.WithSerializable())
	if err != nil {
		return errors.New("etcd get '" + ResourcePath + "' resource failed")
	}

	rs := make(map[string]*ResourceSpec)
	for _, kv := range resp.Kvs {
		r, err := decode(kv.Value)
		if err == nil && r != nil {
			rs[r.Name] = r
		}
	}

	m.logger.Info("leader init load resource", zap.Int("length", len(m.resources)))
	m.resources = rs

	return nil
}

func (m *Master) HandleMsg() {
	msgCh := make(chan *Message)

	select {
	case msg := <-msgCh:
		switch msg.Cmd {
		case MSG_ADD:
			m.AddResource(msg.Specs)
		}
	}
}
