package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	goneGrpc "github.com/gone-io/goner/grpc"
	"github.com/gone-io/goner/nacos"
	"github.com/gone-io/goner/viper"
	"google.golang.org/grpc"
	"grpc_use_discovery/proto"
	"log"
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

	gone.
		NewApp(
			goneGrpc.ServerLoad,
			nacos.LoadRegistry,
			viper.Load,
		).
		Load(&server{}).
		// 启动服务
		Serve()
}
