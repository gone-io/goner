package services

import (
	"context"
	"examples/otel/tracer/oltp/grpc/proto"
	"fmt"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	gone.Flag
	proto.UnimplementedHelloServer              // 嵌入UnimplementedHelloServer
	grpcServer                     *grpc.Server `gone:"*"` // 注入grpc.Server
	logger                         gone.Logger  `gone:"*"`
	tracer                         trace.Tracer
}

const tracerName = "grpc-hello-server"

func (s *server) Init() {
	proto.RegisterHelloServer(s.grpcServer, s) //注册服务
	s.tracer = otel.Tracer(tracerName)
}

// Say  重载协议中定义的服务
func (s *server) Say(ctx context.Context, in *proto.SayRequest) (*proto.SayResponse, error) {
	x, span := s.tracer.Start(ctx, "Say")
	defer span.End()

	s.logger.Infof("Received: %v", in.GetName())
	span.SetAttributes(attribute.Key("Name").String(in.GetName()))

	say := s.doSay(x, in.GetName())

	return &proto.SayResponse{Message: say}, nil
}

func (s *server) doSay(ctx context.Context, name string) string {
	_, span := s.tracer.Start(ctx, "doSay")
	defer span.End()

	span.SetAttributes(attribute.Key("name").String(name))
	span.AddEvent("doSay", trace.WithAttributes(attribute.Key("name").String(name)))

	return fmt.Sprintf("Hello, %s!", name)
}
