package master

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const (
	ResourcePath = "/resources"

	MsgAdd Command = iota
	MSG_DEL
)

type Command int

type Message struct {
	Cmd   Command
	Specs []*ResourceSpec
}

type NodeSpec struct {
	Node    *registry.Node
	Payload int // 负载
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

func getNodeID(assigned string) (string, error) {
	nodes := strings.Split(assigned, "|")
	if len(nodes) < 2 {
		return "", errors.New("assigned node format not collect")
	}

	return nodes[0], nil
}

func (m *Master) reAssign() {
	rs := make([]*ResourceSpec, 0, len(m.resources))

	for _, r := range m.resources {
		if r.AssignedNode == "" {
			continue
		}

		id, err := getNodeID(r.AssignedNode)

		if err != nil {
			m.logger.Error("get node id failed", zap.Error(err))

			continue
		}

		if _, ok := m.workerNodes[id]; !ok {
			rs = append(rs, r)
		}
	}

	m.AddResource(rs)
}

func (m *Master) Assign(r *ResourceSpec) (*NodeSpec, error) {
	candidates := make([]*NodeSpec, 0, len(m.workerNodes))

	for _, n := range m.workerNodes {
		candidates = append(candidates, n)
	}

	// 找到最低的负载
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Payload < candidates[i].Payload
	})

	if len(candidates) > 0 {
		return candidates[0], nil
	}

	return nil, errors.New("no  worker nodes")
}

func (m *Master) AddSeed() {
	rs := make([]*ResourceSpec, 0, len(m.Seeds))
	for _, seed := range m.Seeds {
		resp, err := m.etcdCli.Get(context.Background(), getResourcePath(seed.Name), clientv3.WithPrefix(), clientv3.WithSerializable())
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
		r.ID = m.IDGen.Generate().String()
		ns, err := m.Assign(r)
		if err != nil {
			m.logger.Error("assign resource failed", zap.Error(err))

			continue
		}

		if ns.Node == nil {
			m.logger.Error("no node to assign")

			continue
		}

		r.AssignedNode = ns.Node.Id + "|" + ns.Node.Address
		r.CreateTime = time.Now().Local().UnixNano()

		m.logger.Debug("add resource", zap.Any("specs", r))

		_, err = m.etcdCli.Put(context.Background(), getResourcePath(r.Name), encode(r))
		if err != nil {
			m.logger.Error("put resource to etcd failed", zap.Error(err))

			continue
		}

		m.resources[r.Name] = r
		ns.Payload++
	}
}

func (m *Master) loadResource() error {
	resp, err := m.etcdCli.Get(context.Background(), ResourcePath, clientv3.WithPrefix(), clientv3.WithSerializable())
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

	for _, r := range m.resources {
		if r.AssignedNode != "" {
			id, err := getNodeID(r.AssignedNode)
			if err != nil {
				m.logger.Error("get node id failed", zap.Error(err))

				continue
			}

			node := m.workerNodes[id]
			node.Payload++
		}
	}

	return nil
}

func (m *Master) HandleMsg() {
	msgCh := make(chan *Message)

	select {
	case msg := <-msgCh:
		switch msg.Cmd {
		case MsgAdd:
			m.AddResource(msg.Specs)
		}
	}
}
