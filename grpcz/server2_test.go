package grpcz

import (
	pb "github.com/psoKnight/go-common/grpcz/userpb"
	"testing"
)

func TestServer2(t *testing.T) {
	grpcServerConf := &GrpcServerConf{
		EtcdConfig: &EtcdConf{
			Endpoints: []string{"10.117.49.69:12379", "10.117.49.69:22379", "10.117.49.69:32379"},
		},
		Address:    "127.0.0.1:8001",
		ServerName: "simple_grpc",
	}

	grpcServer, err := NewGrpcServer(grpcServerConf)
	if err != nil {
		t.Errorf("Grpc server init err: %v.", err)
		return
	}

	// register service
	pb.RegisterPlatformServiceServer(grpcServer.GetServer(), &PlatformService{})

	// start grpc server
	if err := grpcServer.Start(); err != nil {
		t.Errorf("Grpc server start err: %v.", err)
		return
	}
}
