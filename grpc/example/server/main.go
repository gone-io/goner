package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	goneGrpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
)

type server struct {
	gone.Flag
	proto.UnimplementedHelloServer // 嵌入UnimplementedHelloServer
}

// 重载协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.SayResponse{Message: "Hello " + in.GetName()}, nil
}

// 实现 goneGrpc.Service接口的RegisterGrpcServer方法，该方法在服务器启动时会被自动调用
func (s *server) RegisterGrpcServer(server *grpc.Server) {
	proto.RegisterHelloServer(server, s)
}

func main() {
	gone.
		Load(&server{}).
		Loads(goner.BaseLoad, goneGrpc.ServerLoad).
		// 启动服务
		Serve()
}
