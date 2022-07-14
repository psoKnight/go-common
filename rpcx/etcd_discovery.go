package rpcx

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

const schema = "grpclb"

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli        *clientv3.Client    // 服务发现客户端
	cc         resolver.ClientConn // gRPC
	serverList sync.Map            // 当前的注册服务
}

// NewServiceDiscovery  新建发现服务
func NewServiceDiscovery(cfg *EtcdConf) resolver.Builder {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout * time.Second,
	})
	if err != nil {
		logrus.Errorf("New service discovery err: %v.", err)
	}

	return &ServiceDiscovery{
		cli: cli,
	}
}

// Build 为给定目标创建一个新的`resolver`，当调用`grpc.Dial()`时执行
func (s *ServiceDiscovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	logrus.Infof("Start build.")

	s.cc = cc

	prefix := "/" + target.Scheme + "/" + target.Endpoint + "/"

	// 根据key 获取对应的键值，此处只返回匹配指定前缀的值
	// 获取当前
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	// 更新grpc 当前的注册服务
	if err := s.cc.UpdateState(resolver.State{Addresses: s.getServices()}); err != nil {
		return nil, err
	}

	// 根据key 获取对应的键值，此处只返回匹配指定前缀的值
	// 获取动态增长
	go s.watcher(prefix)

	return s, nil
}

// ResolveNow 监视目标更新
func (s *ServiceDiscovery) ResolveNow(rn resolver.ResolveNowOptions) {
	logrus.Info("Resolve now.")
}

// Scheme 返回schema
func (s *ServiceDiscovery) Scheme() string {
	return schema
}

// Close 关闭服务
func (s *ServiceDiscovery) Close() {
	s.cli.Close()
}

// watcher 监听前缀
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())

	logrus.Infof("Watching prefix: %s", prefix)

	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: // 新增/修改
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: // 删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// SetServiceList 新增服务地址
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.serverList.Store(key, resolver.Address{Addr: val})

	// 更新grpc 当前的注册服务
	s.cc.UpdateState(resolver.State{Addresses: s.getServices()})
	logrus.Infof("Put key: %s, val: %s.", key, val)
}

// DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.serverList.Delete(key)

	// 更新grpc 当前的注册服务
	s.cc.UpdateState(resolver.State{Addresses: s.getServices()})
	logrus.Infof("Del key: %s.", key)
}

// GetServices 获取服务地址
func (s *ServiceDiscovery) getServices() []resolver.Address {
	//addrs := make([]resolver.Address, 0, 10)
	addrs := make([]resolver.Address, 0)
	s.serverList.Range(func(k, v interface{}) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}
