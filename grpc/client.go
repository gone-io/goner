package gone_grpc

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"reflect"
)

type clientRegister struct {
	gone.Flag
	gone.Logger `gone:"*"`
	clients     []Client      `gone:"*"`
	tracer      tracer.Tracer `gone:"*"`
	connections map[string]*grpc.ClientConn

	requestIdKey string `gone:"config,server.grpc.x-request-id-key=X-Request-Id"`
	tracerIdKey  string `gone:"config,server.grpc.x-trace-id-key=X-Trace-Id"`
}

//go:gone
func NewRegister() gone.Goner {
	return &clientRegister{connections: make(map[string]*grpc.ClientConn)}
}

func (s *clientRegister) traceInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	_ ...grpc.CallOption,
) error {
	ctx = metadata.AppendToOutgoingContext(ctx, s.tracerIdKey, s.tracer.GetTraceId())
	return invoker(ctx, method, req, reply, cc)
}

func (s *clientRegister) register(client Client) error {
	conn, ok := s.connections[client.Address()]
	if !ok {
		c, err := grpc.Dial(
			client.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(s.traceInterceptor),
		)
		if err != nil {
			return err
		}

		s.connections[client.Address()] = c
		conn = c
	}

	client.Stub(conn)
	return nil
}

func (s *clientRegister) Start() error {
	for _, c := range s.clients {
		s.Infof("register gRPC client %v on address %v\n", reflect.ValueOf(c).Type().String(), c.Address())
		if err := s.register(c); err != nil {
			return err
		}
	}

	return nil
}

func (s *clientRegister) Stop() error {
	for _, conn := range s.connections {
		err := conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
