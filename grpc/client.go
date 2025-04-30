package grpc

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)

type clientRegister struct {
	gone.Flag
	logger             gone.Logger          `gone:"*"`
	clients            []Client             `gone:"*"`
	grpcOptions        []grpc.DialOption    `gone:"*"`
	tracer             g.Tracer             `gone:"*" option:"allowNil"`
	discovery          g.ServiceDiscovery   `gone:"*" option:"allowNil"`
	isOtelTracerLoaded g.IsOtelTracerLoaded `gone:"*" option:"allowNil"`

	connections map[string]*grpc.ClientConn
	rb          resolver.Builder

	configure           gone.Configure `gone:"configure"`
	loadBalancingPolicy string         `gone:"config,server.grpc.lb-policy=round_robin"`
	insecure            bool           `gone:"config,server.grpc.insecure=true"`
	requestIdKey        string         `gone:"config,server.grpc.x-request-id-key=X-Request-Id"`
	tracerIdKey         string         `gone:"config,server.grpc.x-trace-id-key=X-Trace-Id"`
}

func (s *clientRegister) Init() {
	if s.discovery != nil {
		s.rb = NewResolverBuilder(s.discovery, s.logger)
	}
}

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

func (s *clientRegister) createConn(address string) (conn *grpc.ClientConn, err error) {
	var options = append(s.grpcOptions, grpc.WithChainUnaryInterceptor(s.traceInterceptor))
	if s.insecure {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if s.isOtelTracerLoaded {
		options = append(options, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	if s.rb != nil {
		options = append(options,
			grpc.WithResolvers(s.rb),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, s.loadBalancingPolicy)),
		)
	}

	return grpc.NewClient(
		address,
		options...,
	)
}

// getConn 根据不同的地址创建 grpc.ClientConn
func (s *clientRegister) getConn(address string) (conn *grpc.ClientConn) {
	conn = s.connections[address]
	if conn == nil {
		var err error
		conn, err = s.createConn(address)
		g.PanicIfErr(gone.ToErrorWithMsg(err, fmt.Sprintf("gRPC createConn for %s", address)))
		s.connections[address] = conn
	}
	return
}

func (s *clientRegister) register(client Client) {
	conn := s.getConn(client.Address())
	client.Stub(conn)
}

// Provide 根据不同的配置创建 grpc.ClientConn
func (s *clientRegister) Provide(tagConf string) (*grpc.ClientConn, error) {
	m, _ := gone.TagStringParse(tagConf)
	address := m["address"]
	if configKey, ok := m["config"]; ok {
		err := s.configure.Get(configKey, &address, address)
		g.PanicIfErr(gone.ToErrorWithMsg(err, "get address from configure err"))
	}
	if address == "" {
		return nil, gone.ToError("address is empty")
	}
	return s.getConn(address), nil
}

func (s *clientRegister) Start() error {
	for _, c := range s.clients {
		s.logger.Infof("register gRPC client %T on address %v", c, c.Address())
		s.register(c)
	}
	return nil
}

func (s *clientRegister) Stop() error {
	for _, conn := range s.connections {
		g.ErrorPrinter(s.logger, conn.Close(), "close gRPC client err")
	}
	return nil
}
