package etcd

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdRegisterConfig struct {
	client        *clientv3.Client                        // 服务注册客户端
	lease         clientv3.Lease                          // 租约
	leaseResp     *clientv3.LeaseGrantResponse            // 租约响应
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 租约响应chan
}

// NewEtcdRegisterService 新建 服务注册
func NewEtcdRegisterService(cfg *EtcdConfig) (*EtcdRegisterConfig, error) {
	conf := &clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
		Username:    cfg.Username,
		Password:    cfg.Password,
	}

	var (
		client *clientv3.Client
	)
	if clientTem, err := clientv3.New(*conf); err == nil {
		client = clientTem
	} else {
		return nil, err
	}

	erc := &EtcdRegisterConfig{
		client: client,
	}
	if err := erc.setLease(cfg.TTL); err != nil {
		return nil, err
	}

	go erc.listenLeaseRespChan()

	return erc, nil
}

// 设置租约
func (cfg *EtcdRegisterConfig) setLease(ttl int64) error {
	lease := clientv3.NewLease(cfg.client)

	leaseResp, err := lease.Grant(context.TODO(), ttl)
	if err != nil {
		return err
	}

	ctx, _ := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	cfg.lease = lease
	cfg.leaseResp = leaseResp
	cfg.keepAliveChan = leaseRespChan
	return nil
}

// 撤销租约
func (cfg *EtcdRegisterConfig) revokeLease() error {
	time.Sleep(time.Duration(2) * time.Second)
	_, err := cfg.lease.Revoke(context.TODO(), cfg.leaseResp.ID)
	return err
}

// 监听 续租情况
func (cfg *EtcdRegisterConfig) listenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-cfg.keepAliveChan:
			if leaseKeepResp == nil {
				logrus.Infof("[etcd]Lease '%d' has been closed.", cfg.leaseResp.ID)
				return
			} else {
				logrus.Infof("[etcd]Lease '%d' success, detail: %+v.", cfg.leaseResp.ID, leaseKeepResp)
			}
		}
	}
}

// PutService 增/改服务中的key
func (cfg *EtcdRegisterConfig) PutService(key, val string) error {
	kv := clientv3.NewKV(cfg.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(cfg.leaseResp.ID))
	return err
}

// DeleteService 删除服务中的key
func (cfg *EtcdRegisterConfig) DeleteService(key string) error {
	kv := clientv3.NewKV(cfg.client)
	_, err := kv.Delete(context.TODO(), key)
	return err
}

// CloceService 关闭服务
func (cfg *EtcdRegisterConfig) CloceService() error {
	_, err := cfg.lease.Revoke(context.TODO(), cfg.leaseResp.ID)
	if err != nil {
		return err
	}

	return cfg.client.Close()
}
