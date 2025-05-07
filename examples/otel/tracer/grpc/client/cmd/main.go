package main

import (
	"context"
	"examples/otel/tracer/oltp/grpc/proto"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
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

	logger gone.Logger `gone:"*"`

	tracer trace.Tracer
}

const tracerName = "grpc-hello-client"

func (c *helloClient) Init() {
	c.HelloClient = proto.NewHelloClient(c.clientConn)
	c.tracer = otel.Tracer(tracerName)
}

func (c *helloClient) Do() {
	ctx, span := c.tracer.Start(context.Background(), "Do Hello Test")
	defer span.End()

	// 调用Say方法，给服务段发送消息
	say, err := c.Say(ctx, &proto.SayRequest{Name: "gone"})
	if err != nil {
		c.logger.Infof("er:%v", err)
		return
	}
	c.logger.Infof("say result: %s", say.Message)
}

func main() {
	gone.Run(func(c *helloClient) {
		c.Do()
	})
}
