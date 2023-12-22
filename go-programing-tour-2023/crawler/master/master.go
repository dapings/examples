package master

import (
	"context"
	"errors"
	"net"
	"time"
	
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"go.uber.org/zap"
)

type Master struct {
	ID string
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

	for {
		select {
		case err := <-leaderCh:
			if err != nil {
				m.logger.Error("leader elect failed", zap.Error(err))

				go m.elect(e, leaderCh)
			} else {
				m.logger.Info("master change to leader")
			}
		case resp := <-leaderChange:
			if len(resp.Kvs) > 0 {
				m.logger.Info("watch leader change", zap.String("leader:", string(resp.Kvs[0].Value)))
			}
		case <-time.After(10 * time.Second):
			resp, err := e.Leader(context.Background())
			if err != nil {
				m.logger.Error("get leader failed", zap.Error(err))
			}
			if resp != nil && len(resp.Kvs) > 0 {
				m.logger.Info("get leader", zap.String("value", string(resp.Kvs[0].Value)))
			}
		}
	}
}

func (m *Master) elect(e *concurrency.Election, ch chan error) {
	// block util election success.
	err := e.Campaign(context.Background(), m.ID)
	ch <- err
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
