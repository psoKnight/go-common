package grpcz

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type GrpcClientConf struct {
	EtcdConfig *EtcdConf `json:"etcd_config"` // etcd 相关配置
	ServerName string    `json:"server_name"` // 服务名称
}

type GrpcClient struct {
	cfg        *GrpcClientConf  // client配置
	grpcClient *grpc.ClientConn // gRPC client 端
}

// NewGrpcClient 新建gRPC client
func NewGrpcClient(cfg *GrpcClientConf) (*GrpcClient, error) {

	// etcd 服务发现
	d, err := NewServiceDiscovery(cfg.EtcdConfig)
	if err != nil {
		return nil, err
	}

	resolver.Register(d)

	// 连接服务器
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", d.Scheme(), cfg.ServerName),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 轮询策略
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &GrpcClient{cfg: cfg, grpcClient: conn}, nil
}

// GetClient 获取client
func (s *GrpcClient) GetClient() *grpc.ClientConn {
	return s.grpcClient
}
