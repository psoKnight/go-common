package rpcx

import (
	pb "github.com/psoKnight/go-common/rpcx/userpb"
	"testing"
)

func Test_Server2(t *testing.T) {
	grpcServerConf := &RpcxServerConf{
		EtcdConfig: &EtcdConf{
			Endpoints: []string{"10.117.49.69:12379", "10.117.49.69:22379", "10.117.49.69:32379"},
		},
		Address:    "127.0.0.1:8001",
		ServerName: "simple_grpc",
	}

	grpcServer, err := NewRpcxServer(grpcServerConf)
	if err != nil {
		t.Errorf("rpcx server init err: %v.", err)
		return
	}

	// register service
	pb.RegisterPlatformServiceServer(grpcServer.GrpcServer(), &PlatformService{})

	// start grpc server
	if err := grpcServer.Start(); err != nil {
		t.Errorf("Rpcx server start err: %v.", err)
		return
	}
}
