package main

import (
	"context"
	"example/grpc/proto"
	"fmt"
	gone_grpc "github.com/gone-io/goner/grpc"
	"github.com/gone-io/goner/viper"
	"google.golang.org/grpc"
	"log"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // 嵌入HelloClient

	host string `gone:"config,server.host"`
	port string `gone:"config,server.port"`
}

// 实现 gone_grpc.Client接口的Address方法，该方法在客户端启动时会被自动调用
// 该方法的作用是告诉客户端gRPC服务的地址
func (c *helloClient) Address() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

// 实现 gone_grpc.Client接口的Stub方法，该方法在客户端启动时会被自动调用
// 在该方法中，完成 HelloClient的初始化
func (c *helloClient) Stub(conn *grpc.ClientConn) {
	c.HelloClient = proto.NewHelloClient(conn)
}

func main() {
	gone.
		Load(&helloClient{}).
		Loads(viper.Load, gone_grpc.ClientRegisterLoad).
		Run(func(in struct {
			hello *helloClient `gone:"*"` // 在Run方法的参数中，注入 helloClient
		}) {
			// 调用Say方法，给服务段发送消息
			say, err := in.hello.Say(context.Background(), &proto.SayRequest{Name: "gone"})
			if err != nil {
				log.Printf("er:%v", err)
				return
			}
			log.Printf("say result: %s", say.Message)
		})
}
