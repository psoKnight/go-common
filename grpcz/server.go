package grpcz

import (
	"errors"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerConf struct {
	EtcdConfig *EtcdConf `json:"etcd_config"` // etcd_ 相关配置
	Address    string    `json:"address"`     // 监听地址
	ServerName string    `json:"server_name"` // 服务名称
}

func (c *GrpcServerConf) Check() error {
	if c.Address == "" {
		return errors.New("miss server addr")
	}
	if c.ServerName == "" {
		return errors.New("miss server name")
	}
	return nil
}

type GrpcServer struct {
	cfg        *GrpcServerConf // server配置
	grpcServer *grpc.Server    // gRPC server 端
}

// NewGrpcServer 新建gRPC server
func NewGrpcServer(cfg *GrpcServerConf) (*GrpcServer, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}

	var opts []grpc.ServerOption

	// grpc-middleware
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(), //recover
	)))
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_recovery.StreamServerInterceptor(), //recover
	)))

	//TODO auth, timeout

	// 新建gRPC 服务器实例
	grpcServer := grpc.NewServer(opts...)

	grpcCli := &GrpcServer{cfg: cfg, grpcServer: grpcServer}

	// 服务注册到etcd
	_, err := NewServiceRegister(cfg.EtcdConfig, cfg.ServerName, cfg.Address)
	if err != nil {
		return nil, err
	}

	return grpcCli, nil
}

// Start 运行gRPC 服务
func (s *GrpcServer) Start() error {
	// 监听本地端口
	listen, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return err
	}
	logrus.Infof("[grpc]%s net.Listenning...", s.cfg.Address)

	// 用服务器Serve() 方法以及端口信息区实现阻塞等待，直到进程被杀死或者Stop() 被调用
	if err := s.grpcServer.Serve(listen); err != nil {
		logrus.Errorf("[grpc]server err: %v.", err)
		return err
	}

	return nil
}

// GetServer 返回*grpc.Server
func (s *GrpcServer) GetServer() *grpc.Server {
	return s.grpcServer
}
