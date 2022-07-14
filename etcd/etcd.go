package etcd

import (
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdConfig struct {
	Endpoints []string `json:"endpoints"`

	// DialTimeout 是建立连接失败的超时时间
	DialTimeout time.Duration `json:"dial_timeout"`

	// 用户名是用于认证的用户名
	Username string `json:"username"`

	// 密码是用于认证的密码
	Password string `json:"password"`

	// TODO: support custom balancer picker
	TTL int64
}

type Etcd struct {
	register *EtcdRegisterConfig // 服务注册
	discover *EtcdDiscoverConfig // 服务发现
	cfg      *EtcdConfig         // 基础配置
}

// NewEtcd 新建etcd
func NewEtcd(cfg *EtcdConfig) (*Etcd, error) {

	if cfg == nil {
		return nil, errors.New("[etcd]config is nil")
	}

	client := &Etcd{
		cfg: cfg,
	}

	newEtcdRegisterService, err := NewEtcdRegisterService(cfg)
	if err != nil {
		return nil, err
	}
	client.register = newEtcdRegisterService

	newEtcdDiscoverService, err := NewEtcdDiscoverService(cfg)
	if err != nil {
		return nil, err
	}
	client.discover = newEtcdDiscoverService

	return client, nil
}

// PutRegisterService 增/改注册服务中的key
func (e *Etcd) PutRegisterService(key, val string) error {
	return e.register.PutService(key, val)
}

// GetRegister 获取服务注册client
func (e *Etcd) GetRegister() *clientv3.Client {
	return e.register.client
}

// GetDiscover 获取服务发现client
func (e *Etcd) GetDiscover() *clientv3.Client {
	return e.discover.client
}

// DeleteRegisterService 删除注册服务中的key
func (e *Etcd) DeleteRegisterService(key string) error {
	return e.register.DeleteService(key)
}

// CloceRegisterService 关闭注册服务
func (e *Etcd) CloceRegisterService() error {
	return e.register.CloceService()
}

// WatchDiscoverService 监听服务列表
func (e *Etcd) WatchDiscoverService(prefix string) (map[string]string, error) {
	return e.discover.WatchService(prefix)
}

// GetDiscoverService 获取服务列表
func (e *Etcd) GetDiscoverService(prefix string) (map[string]string, error) {
	return e.discover.GetService(prefix)
}

// CloseDiscoverService 关闭服务发现
func (e *Etcd) CloseDiscoverService() error {
	return e.discover.CloseService()
}
