package etcd

import (
	"errors"
	"github.com/sirupsen/logrus"
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

type EtcdClient struct {
	register *EtcdRegisterConfig // 服务注册
	discover *EtcdDiscoverConfig // 服务发现
	cfg      *EtcdConfig         // 基础配置
}

// 新建etcd 客户端
func NewEtcdClient(cfg *EtcdConfig) (*EtcdClient, error) {

	if cfg == nil {
		return nil, errors.New("[etcd]Config is nil.")
	}

	client := &EtcdClient{
		cfg: cfg,
	}

	newEtcdRegisterService, err := NewEtcdRegisterService(cfg)
	if err != nil {
		logrus.Errorf("[etcd]New register service err: %v.", err)
		return nil, err
	}
	client.register = newEtcdRegisterService

	newEtcdDiscoverService, err := NewEtcdDiscoverService(cfg)
	if err != nil {
		logrus.Errorf("[etcd]New discover service err: %v.", err)
		return nil, err
	}
	client.discover = newEtcdDiscoverService

	logrus.Infof("[etcd]%v connect success.", cfg.Endpoints)

	return client, nil
}

// PutRegisterService 增/改注册服务中的key
func (cli *EtcdClient) PutRegisterService(key, val string) error {
	return cli.register.PutService(key, val)
}

// DeleteRegisterKey 删除注册服务中的key
func (cli *EtcdClient) DeleteRegisterService(key string) error {
	return cli.register.DeleteService(key)
}

// CloceRegisterService 关闭注册服务
func (cli *EtcdClient) CloceRegisterService() error {
	return cli.register.CloceService()
}

// WatchService 监听服务列表
func (cli *EtcdClient) WatchDiscoverService(prefix string) (map[string]string, error) {
	return cli.discover.WatchService(prefix)
}

// GetDiscoverService 获取服务列表
func (cli *EtcdClient) GetDiscoverService(prefix string) (map[string]string, error) {
	return cli.discover.GetService(prefix)
}

// CloseDiscoverService 关闭服务发现
func (cli *EtcdClient) CloseDiscoverService() error {
	return cli.discover.CloseService()
}
