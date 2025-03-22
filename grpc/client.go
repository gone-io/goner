package grpc

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"reflect"
)

type clientRegister struct {
	gone.Flag
	logger    gone.Logger    `gone:"*"`
	clients   []Client       `gone:"*"`
	tracer    g.Tracer       `gone:"*" option:"allowNil"`
	configure gone.Configure `gone:"configure"`

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
	tracerId, _ := ctx.Value(s.tracerIdKey).(string)
	if s.tracer != nil {
		tracerId = s.tracer.GetTraceId()
	}
	ctx = metadata.AppendToOutgoingContext(ctx, s.tracerIdKey, tracerId)
	return invoker(ctx, method, req, reply, cc)
}

// getConn 根据不同的地址创建 grpc.ClientConn
func (s *clientRegister) getConn(address string) (conn *grpc.ClientConn, err error) {
	conn = s.connections[address]
	if conn == nil {
		if conn, err = grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(s.traceInterceptor),
		); err != nil {
			return nil, gone.ToError(err)
		}
		s.connections[address] = conn
	}
	return
}

func (s *clientRegister) register(client Client) error {
	conn, err := s.getConn(client.Address())
	if err != nil {
		return err
	}
	client.Stub(conn)
	return nil
}

// Provide 根据不同的配置创建 grpc.ClientConn
func (s *clientRegister) Provide(tagConf string) (*grpc.ClientConn, error) {
	m, _ := gone.TagStringParse(tagConf)
	address := m["address"]
	if configKey, ok := m["config"]; ok {
		err := s.configure.Get(configKey, &address, address)
		if err != nil {
			return nil, gone.ToError(err)
		}
	}
	if address == "" {
		return nil, gone.ToError("address is empty")
	}
	return s.getConn(address)
}

func (s *clientRegister) Start() error {
	for _, c := range s.clients {
		s.logger.Infof("register gRPC client %v on address %v\n", reflect.ValueOf(c).Type().String(), c.Address())
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
