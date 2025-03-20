package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	gone_grpc "github.com/gone-io/goner/grpc"
	"google.golang.org/grpc"
	"grpc_demo/proto"
	"log"
	"os"
)

type helloClient struct {
	gone.Flag
	proto.HelloClient // 嵌入HelloClient

	clientConn *grpc.ClientConn `gone:"*,config=grpc.service.hello.address"`

	//clientConn2 *grpc.ClientConn `gone:"*,address=127.0.0.1:9090"`
}

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
}

func main() {
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
