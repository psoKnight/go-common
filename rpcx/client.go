package rpcx

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type RpcxClientConf struct {
	EtcdConfig *EtcdConf `json:"etcd_config"` // etcd_ 相关配置
	ServerName string    `json:"server_name"` // 服务名称
}

type RpcxClient struct {
	cfg        *RpcxClientConf  // client配置
	clientConn *grpc.ClientConn // gRPC client 端
}

// NewRpcxClient 新建gRPC client
func NewRpcxClient(cfg *RpcxClientConf) (*RpcxClient, error) {

	// etcd_ 服务发现
	d := NewServiceDiscovery(cfg.EtcdConfig)

	resolver.Register(d)

	// 连接服务器
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", d.Scheme(), cfg.ServerName),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 轮询策略
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		logrus.Errorf("Grpc dial err: %v.", err)
		return nil, err
	}

	return &RpcxClient{cfg: cfg, clientConn: conn}, nil
}

// ClientConn 获取client
func (s *RpcxClient) ClientConn() *grpc.ClientConn {
	return s.clientConn
}
