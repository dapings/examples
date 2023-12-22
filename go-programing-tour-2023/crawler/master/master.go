package master

import (
	"context"
	"errors"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

type Master struct {
	ID          string
	ready       int32
	leaderID    string
	workerNodes map[string]*registry.Node
	options
}

func New(id string, opts ...Option) (*Master, error) {
	m := &Master{}

	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	m.options = options

	ipv4, err := getLocalIP()
	if err != nil {
		return nil, err
	}

	m.ID = genMasterID(id, ipv4, m.GRPCAddr)
	m.logger.Sugar().Debugln("master_id:", m.ID)

	go m.Campaign()

	return &Master{}, nil
}

func (m *Master) Campaign() {
	endpoints := []string{m.registryURL}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		m.logger.Error("etcd v3 client new failed", zap.Error(err))

		panic(err)
	}

	s, sessionErr := concurrency.NewSession(cli, concurrency.WithTTL(5))
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
	e := concurrency.NewElection(s, "resources/election")
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
				m.logger.Info("master change to leader")

				m.leaderID = m.ID
				if !m.IsLeader() {
					m.BecomeLeader()
				}
			}
		case resp := <-leaderChange:
			if len(resp.Kvs) > 0 {
				m.logger.Info("watch leader change", zap.String("leader:", string(resp.Kvs[0].Value)))
			}
		case resp := <-workerNodeChange:
			m.logger.Info("watch worker change", zap.Any("worker:", resp))

			m.updateNodes()
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

func (m *Master) BecomeLeader() {
	atomic.StoreInt32(&m.ready, 1)
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

func (m *Master) updateNodes() {
	services, err := m.registry.GetService(cmd.WorkerServiceName)
	if err != nil {
		m.logger.Error("get service", zap.Error(err))

		return
	}

	nodes := make(map[string]*registry.Node)
	if len(services) > 0 {
		for _, spec := range services[0].Nodes {
			nodes[spec.Id] = spec
		}
	}

	added, deleted, changed := diffWorkerNode(m.workerNodes, nodes)
	m.logger.Sugar().Info("worker joined: ", added, ", leaved: ", deleted, ", changed: ", changed)

	m.workerNodes = nodes
}

func diffWorkerNode(old, new map[string]*registry.Node) ([]string, []string, []string) {
	added := make([]string, 0)
	deleted := make([]string, 0)
	changed := make([]string, 0)

	for k, v := range new {
		if ov, ok := old[k]; ok {
			if !reflect.DeepEqual(v, ov) {
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
