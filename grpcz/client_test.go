package grpcz

import (
	"context"
	pb "github.com/psoKnight/go-common/grpcz/userpb"
	"github.com/sirupsen/logrus"
	"strconv"
	"testing"
)

func TestClient(t *testing.T) {
	rpcClientConf := &GrpcClientConf{
		EtcdConfig: &EtcdConf{
			Endpoints:   []string{"10.117.49.69:12379", "10.117.49.69:22379", "10.117.49.69:32379"},
			DialTimeout: 5,
			Username:    "",
			Password:    "",
			TTL:         9,
		},
		ServerName: "simple_grpc",
	}

	grpcClient, err := NewGrpcClient(rpcClientConf)
	if err != nil {
		t.Errorf("Grpc client conn err: %v.", err)
		return
	}

	// 建立gRPC 连接
	cli := pb.NewPlatformServiceClient(grpcClient.GetClient())

	for i := 0; i < 10001; i++ {
		itoa := strconv.Itoa(i)
		route(itoa, cli)
	}

}

// route 调用服务端Route 方法
func route(s string, cli pb.PlatformServiceClient) {
	// 创建发送结构体
	req := pb.PlatformServiceRequest{
		Key: s,
	}

	// 调用服务(Route方法)
	// 同时传入了一个context.Context，在有需要时可以让我们改变RPC 的行为，比如超时/取消一个正在运行的RPC
	res, err := cli.Route(context.Background(), &req)
	if err != nil {
		logrus.Errorf("Call Route err: %v.", err)
	}

	// 打印返回值
	logrus.Info(res)
}
