package grpcz

import (
	"context"
	pb "github.com/psoKnight/go-common/grpcz/userpb"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestServer(t *testing.T) {
	grpcServerConf := &GrpcServerConf{
		EtcdConfig: &EtcdConf{
			Endpoints: []string{"10.117.49.69:12379", "10.117.49.69:22379", "10.117.49.69:32379"},
		},
		Address:    "127.0.0.1:8000",
		ServerName: "simple_grpc",
	}

	grpcServer, err := NewGrpcServer(grpcServerConf)
	if err != nil {
		t.Errorf("rpcx server init err: %v.", err)
		return
	}

	// register service
	pb.RegisterPlatformServiceServer(grpcServer.GetServer(), &PlatformService{})

	// start grpc server
	if err := grpcServer.Start(); err != nil {
		t.Errorf("Rpcx server start err: %v.", err)
		return
	}
}

// PlatformService 服务
type PlatformService struct{}

// Route 实现Route方法
func (s *PlatformService) Route(ctx context.Context, req *pb.PlatformServiceRequest) (*pb.PlatformServiceResponse, error) {
	logrus.Infof("Receive: %s.", req.Key)
	res := pb.PlatformServiceResponse{
		Code:  200,
		Value: req.Key,
	}
	return &res, nil
}
