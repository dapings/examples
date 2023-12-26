package master

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	pb "github.com/dapings/examples/go-programing-tour-2023/crawler/protos/crawler"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/client"
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
		return "", errors.New("assigned node format not incorrect")
	}

	return nodes[0], nil
}

func (m *Master) reAssign() {
	rs := make([]*ResourceSpec, 0, len(m.resources))

	for _, r := range m.resources {
		if r.AssignedNode == "" {
			rs = append(rs, r)

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

	m.AddResources(rs)
}

func (m *Master) Assign(_ *ResourceSpec) (*NodeSpec, error) {
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

	m.AddResources(rs)
}

func (m *Master) AddResources(rs []*ResourceSpec) {
	for _, r := range rs {
		_, err := m.addResource(r)
		if err != nil {
			continue
		}
	}
}

func (m *Master) AddResource(ctx context.Context, req *pb.ResourceSpec, resp *pb.NodeSpec) error {
	if !m.IsLeader() && m.leaderID != "" && m.leaderID != m.ID {
		// 当前已不再是leader，转发到 leader 上。
		addr := getLeaderAddr(m.leaderID)
		nodeSpec, err := m.forwardCli.AddResource(ctx, req, client.WithAddress(addr))
		if nodeSpec != nil {
			resp.Id = nodeSpec.Id
			resp.Address = nodeSpec.Address
		}

		return err
	}

	nodeSpec, err := m.addResource(&ResourceSpec{Name: req.Name})
	if nodeSpec != nil {
		resp.Id = nodeSpec.Node.Id
		resp.Address = nodeSpec.Node.Address
	}

	return err
}

func (m *Master) DelResource(_ context.Context, spec *pb.ResourceSpec, _ *empty.Empty) error {
	r, ok := m.resources[spec.Name]
	if !ok {
		return errors.New("no such task")
	}

	if _, err := m.etcdCli.Delete(context.Background(), getResourcePath(spec.Name)); err != nil {
		return err
	}

	if r.AssignedNode != "" {
		nodeID, err := getNodeID(r.AssignedNode)
		if err != nil {
			return err
		}

		if ns, ok := m.workerNodes[nodeID]; ok {
			ns.Payload -= 1
		}
	}

	return nil
}

func (m *Master) addResource(r *ResourceSpec) (*NodeSpec, error) {
	r.ID = m.IDGen.Generate().String()
	ns, err := m.Assign(r)
	if err != nil {
		m.logger.Error("assign resource failed", zap.Error(err))

		return nil, err
	}

	if ns.Node == nil {
		m.logger.Error("no node to assign")

		return nil, err
	}

	r.AssignedNode = ns.Node.Id + "|" + ns.Node.Address
	r.CreateTime = time.Now().Local().UnixNano()

	m.logger.Debug("add resource", zap.Any("specs", r))

	_, err = m.etcdCli.Put(context.Background(), getResourcePath(r.Name), encode(r))
	if err != nil {
		m.logger.Error("put resource to etcd failed", zap.Error(err))

		return nil, err
	}

	m.resources[r.Name] = r
	ns.Payload++

	return ns, nil
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

			if node, ok := m.workerNodes[id]; ok {
				node.Payload++
			}
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
			m.AddResources(msg.Specs)
		}
	}
}
