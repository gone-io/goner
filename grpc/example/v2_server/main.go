package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	goneGrpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
	"os"
)

type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // 嵌入UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // 注入grpc.Server
}

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) //注册服务
}

// Say  重载协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// gone内置默认的配置组件只能从环境变量中读取配置，所以需要设置环境变量
	os.Setenv("GONE_SERVER_GRPC_PORT", "9091")

	gone.
		Load(&server{}).
		Loads(goneGrpc.ServerLoad).
		// 启动服务
		Serve()
}
