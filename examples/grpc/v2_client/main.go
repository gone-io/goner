package main

import (
	"context"
	"example/grpc/proto"
	gone_grpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"log"
	"os"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // 使用方法1：嵌入HelloClient，本组件只负载初始化，能力提供给第三方组件使用

	// 使用方法2：在本组件直接使用，不提供第三方组件使用
	//hello *proto.HelloClient

	// config=${配置的key},address=${服务地址}； //config优先级更高
	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`

	// config 和 address 可以一起使用，如果config没有读取到值，降级为使用 address
	//clientConn1 *grpc.ClientConn `gone:"*,config=grpc.service.hello.address,address=127.0.0.1:9091""`

	// address 也可以单独使用，不推荐这种方式，意味着写死了
	//clientConn2 *grpc.ClientConn `gone:"*,address=127.0.0.1:9091"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
	//c.hello = &c.HelloClient
}

func main() {
	// gone内置默认的配置组件只能从环境变量中读取配置，所以需要设置环境变量
	os.Setenv("GONE_GRPC_SERVICE_HELLO_ADDRESS", "127.0.0.1:9091")

	gone.
		Load(&helloClient{}).
		Loads(gone_grpc.ClientRegisterLoad).
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
