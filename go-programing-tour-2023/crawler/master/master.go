package master

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd"
	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

type Master struct {
	ID          string
	ready       int32
	leaderID    string
	workerNodes map[string]*NodeSpec
	resources   map[string]*ResourceSpec
	IDGen       *snowflake.Node
	etcdCli     *clientv3.Client
	options
}

func New(id string, opts ...Option) (*Master, error) {
	m := &Master{}

	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	m.options = options
	m.resources = make(map[string]*ResourceSpec)

	// ID gen by snowflake.
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	m.IDGen = node

	ipv4, err := getLocalIP()
	if err != nil {
		return nil, err
	}

	m.ID = genMasterID(id, ipv4, m.GRPCAddr)
	m.logger.Sugar().Debugln("master_id:", m.ID)

	// etcd cli
	endpoints := []string{m.registryURL}
	cli, cliErr := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if cliErr != nil {
		m.logger.Error("etcd v3 client new failed", zap.Error(cliErr))

		return nil, cliErr
	}
	m.etcdCli = cli

	m.updateWorkerNodes()
	m.AddSeed()

	go m.Campaign()
	go m.HandleMsg()

	return m, nil
}

func (m *Master) Campaign() {
	s, sessionErr := concurrency.NewSession(m.etcdCli, concurrency.WithTTL(5))
	if sessionErr != nil {
		m.logger.Error("etcd v3 concurrency new session failed", zap.Error(sessionErr))

		panic(sessionErr)
	}

	defer func(s *concurrency.Session) {
		err := s.Close()
		if err != nil {
			_ = s.Close()
		}
	}(s)

	// 创建一个新的etcd选举
	e := concurrency.NewElection(s, "/crawler/election")
	leaderCh := make(chan error)

	go m.elect(e, leaderCh)

	leaderChange := e.Observe(context.Background())

	select {
	case resp := <-leaderChange:
		m.logger.Info("watch leader change", zap.String("leader:", string(resp.Kvs[0].Value)))
	}

	workerNodeChange := m.WatcherWorker()

	for {
		select {
		case err := <-leaderCh:
			if err != nil {
				m.logger.Error("leader elect failed", zap.Error(err))

				go m.elect(e, leaderCh)
			} else {
				m.logger.Info("master start change to leader")

				m.leaderID = m.ID
				if !m.IsLeader() {
					if err := m.BecomeLeader(); err != nil {
						m.logger.Error("become leader failed", zap.Error(err))

						// NOTE: 切换到 leader 失败后，如何处理？是再次选主？
					}
				}
			}
		case resp := <-leaderChange:
			if len(resp.Kvs) > 0 {
				m.logger.Info("watch leader change", zap.String("leader:", string(resp.Kvs[0].Value)))
			}
		case resp := <-workerNodeChange:
			m.logger.Info("watch worker change", zap.Any("worker:", resp))

			m.updateWorkerNodes()
			m.reAssign()
		case <-time.After(20 * time.Second):
			resp, err := e.Leader(context.Background())

			if err != nil {
				m.logger.Error("get leader failed", zap.Error(err))

				if errors.Is(err, concurrency.ErrElectionNoLeader) {
					go m.elect(e, leaderCh)
				}
			}

			if resp != nil && len(resp.Kvs) > 0 {
				m.logger.Info("get leader", zap.String("value", string(resp.Kvs[0].Value)))

				if m.IsLeader() && m.ID != string(resp.Kvs[0].Value) {
					// 当前已不再是 leader
					atomic.StoreInt32(&m.ready, 0)
				}
			}
		}
	}
}

func (m *Master) elect(e *concurrency.Election, ch chan error) {
	// block util election success.
	err := e.Campaign(context.Background(), m.ID)
	ch <- err
}

func (m *Master) IsLeader() bool {
	return atomic.LoadInt32(&m.ready) != 0
}

func (m *Master) BecomeLeader() error {
	m.updateWorkerNodes()

	if err := m.loadResource(); err != nil {
		return fmt.Errorf("load resource failed:%w", err)
	}

	m.reAssign()

	atomic.StoreInt32(&m.ready, 1)
	return nil
}

func (m *Master) WatcherWorker() chan *registry.Result {
	watch, err := m.registry.Watch(registry.WatchService(cmd.WorkerServiceName))
	if err != nil {
		m.logger.Error("registry.Watch service failed", zap.Error(err))

		panic(err)
	}

	ch := make(chan *registry.Result)

	go func() {
		for {
			result, err := watch.Next()
			if err != nil {
				m.logger.Error("watch worker service failed", zap.Error(err))

				continue
			}

			ch <- result
		}
	}()

	return ch
}

func (m *Master) updateWorkerNodes() {
	services, err := m.registry.GetService(cmd.WorkerServiceName)
	if err != nil {
		m.logger.Error("get service", zap.Error(err))

		return
	}

	nodes := make(map[string]*NodeSpec)
	if len(services) > 0 {
		for _, spec := range services[0].Nodes {
			nodes[spec.Id] = &NodeSpec{Node: spec}
		}
	}

	added, deleted, changed := diffWorkerNode(m.workerNodes, nodes)
	m.logger.Sugar().Info("worker joined: ", added, ", leaved: ", deleted, ", changed: ", changed)

	m.workerNodes = nodes
}

func diffWorkerNode(old, new map[string]*NodeSpec) ([]string, []string, []string) {
	added := make([]string, 0)
	deleted := make([]string, 0)
	changed := make([]string, 0)

	for k, v := range new {
		if ov, ok := old[k]; ok {
			if !reflect.DeepEqual(v.Node, ov.Node) {
				changed = append(changed, k)
			}
		} else {
			added = append(added, k)
		}
	}

	for k := range old {
		if _, ok := new[k]; !ok {
			deleted = append(deleted, k)
		}
	}

	return added, deleted, changed
}

func genMasterID(id, ipv4, gRPCAddr string) string {
	return "master" + id + "-" + ipv4 + gRPCAddr
}

// 获取本机网卡IP。
func getLocalIP() (string, error) {
	var (
		address []net.Addr
		err     error
	)

	// 获取所有网卡
	if address, err = net.InterfaceAddrs(); err != nil {
		return "", err
	}

	// 取第一个非lo的网卡IP
	for _, addr := range address {
		if ipNet, isIPNet := addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no local ip")
}
