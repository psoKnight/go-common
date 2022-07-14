package grpcz

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdConf struct {
	Endpoints   []string      `json:"endpoints"`    // etcd_ 主机地址
	DialTimeout time.Duration `json:"dial_timeout"` // 连接失败超时时间
	Username    string        `json:"username"`     // 用户名
	Password    string        `json:"password"`     // 密码
	TTL         int64         `json:"ttl"`          // Lease TTL时间，单位：s；每次KeepAlive 续租频率为TTL/3
}

// Check 检查config 并且设置默认值
func (c *EtcdConf) Check() error {
	if c.DialTimeout == 0 {
		c.DialTimeout = 5
	}

	if len(c.Endpoints) == 0 {
		return errors.New("[grpc]etcd miss endpoints.")
	}

	if c.TTL == 0 {
		c.TTL = 9
	}

	return nil
}

// ServiceRegister 服务注册
type ServiceRegister struct {
	cli           *clientv3.Client                        // 服务注册客户端
	leaseID       clientv3.LeaseID                        // 租约ID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 租约keepalieve 相应chan
	key           string                                  // key
	val           string                                  // value
}

// NewServiceRegister 新建注册服务
func NewServiceRegister(cfg *EtcdConf, serName, addr string) (*ServiceRegister, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout * time.Second,
	})
	if err != nil {
		return nil, err
	}

	ser := &ServiceRegister{
		cli: cli,
		key: "/" + schema + "/" + serName + "/" + addr,
		val: addr,
	}

	// 申请租约设置时间keepalive
	if err := ser.putKeyWithLease(cfg.TTL); err != nil {
		return nil, err
	}

	logrus.Infof("[grpc]new service register success, endpoints: %v, serName: %s. ttl: %d.", cfg.Endpoints, serName, cfg.TTL)
	return ser, nil
}

// 设置租约
func (s *ServiceRegister) putKeyWithLease(ttl int64) error {
	// 设置租约时间
	resp, err := s.cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 设置续租 定期发送需求请求
	leaseRespChan, errSCK := s.cli.KeepAlive(context.Background(), resp.ID)
	if errSCK != nil {
		return errSCK
	}

	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan

	logrus.Infof("[grpc]put key: %s, val: %s success!", s.key, s.val)

	return nil
}

// ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-s.keepAliveChan:
			if leaseKeepResp == nil {
				logrus.Infof("[grpc]lease '%d' has been closed.", s.leaseID)
				return
			} else {
				logrus.Infof("[grpc]lease '%d' success, detail: %+v.", s.leaseID, leaseKeepResp)
			}
		}
	}
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	// 撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	logrus.Info("[grpc]revoke lease.")

	return s.cli.Close()
}
