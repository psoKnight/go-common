package etcd

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type EtcdDiscoverConfig struct {
	client     *clientv3.Client
	serverList map[string]string // 当前的注册服务
	lock       sync.Mutex
}

// NewEtcdDiscoverService 新建 服务发现
func NewEtcdDiscoverService(cfg *EtcdConfig) (*EtcdDiscoverConfig, error) {
	conf := &clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
		Username:    cfg.Username,
		Password:    cfg.Password,
	}
	if client, err := clientv3.New(*conf); err == nil {
		return &EtcdDiscoverConfig{
			client:     client,
			serverList: make(map[string]string),
		}, nil
	} else {
		return nil, err
	}
}

// WatchService 监听服务
func (cfg *EtcdDiscoverConfig) WatchService(prefix string) (map[string]string, error) {

	// 根据key 获取对应的键值，此处只返回匹配指定前缀的值
	// 获取当前
	resp, err := cfg.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	// 根据key 获取对应的键值，此处只返回匹配指定前缀的值
	// 获取动态增长
	addrs := cfg.extractAddrs(resp)

	go cfg.watcher(prefix)
	return addrs, nil
}

// 转换key/value 存储格式
func (cfg *EtcdDiscoverConfig) extractAddrs(resp *clientv3.GetResponse) map[string]string {

	if resp == nil || resp.Kvs == nil {
		return nil
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			cfg.setServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
		}
	}

	return cfg.serverList
}

//watcher 监听前缀
func (cfg *EtcdDiscoverConfig) watcher(prefix string) {
	rch := cfg.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	logrus.Infof("[etcd]Watching prefix: %s.", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: // 新增/修改
				cfg.setServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: // 删除
				cfg.delServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// SetServiceList 根据key 新增/修改当前的服务
func (cfg *EtcdDiscoverConfig) setServiceList(key, val string) {
	cfg.lock.Lock()
	defer cfg.lock.Unlock()
	cfg.serverList[key] = val
	logrus.Infof(fmt.Sprintf("[etcd]Set data key: %s, val: %s.", key, val))
}

// DelServiceList 根据key 删除当前的服务
func (cfg *EtcdDiscoverConfig) delServiceList(key string) {
	cfg.lock.Lock()
	defer cfg.lock.Unlock()
	delete(cfg.serverList, key)
	logrus.Infof(fmt.Sprintf("[etcd]Del data key: %s.", key))
}

// GetService 获取服务
func (cfg *EtcdDiscoverConfig) GetService(prefix string) (map[string]string, error) {
	// TODO
	/**
	1、确认返回格式
	2、是否支持批量获取多类服务
	3.是否全量返回
	*/

	if prefix == "*" {
		return cfg.serverList, nil
	} else {
		// TODO
		return cfg.serverList, nil
	}
}

// Close 关闭服务
func (cfg *EtcdDiscoverConfig) CloseService() error {
	return cfg.client.Close()
}
