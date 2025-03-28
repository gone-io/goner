package grpc

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"reflect"
)

//const XTraceId = "x-trace-id"

func createListener(s *server) (err error) {
	s.listener, err = net.Listen("tcp", s.address)
	return
}

type server struct {
	gone.Flag
	logger       gone.Logger `gone:"*"`
	grpcServices []Service   `gone:"*"`
	cMuxServer   g.Cmux      `gone:"*" option:"allowNil"`
	tracer       g.Tracer    `gone:"*" option:"allowNil"`

	port         int    `gone:"config,server.grpc.port,default=9090"`
	host         string `gone:"config,server.grpc.host,default=0.0.0.0"`
	requestIdKey string `gone:"config,server.grpc.x-request-id-key=X-Request-Id"`
	tracerIdKey  string `gone:"config,server.grpc.x-trace-id-key=X-Trace-Id"`

	grpcServer     *grpc.Server
	listener       net.Listener
	address        string
	createListener func(*server) error
}

func (s *server) GonerName() string {
	return "gone-grpc-server"
}

func (s *server) initListener() error {
	if s.cMuxServer != nil {
		s.listener = s.cMuxServer.MatchFor(g.GRPC)
		s.address = s.cMuxServer.GetAddress()
		return nil
	}

	s.address = fmt.Sprintf("%s:%d", s.host, s.port)
	return s.createListener(s)
}
func (s *server) Init() error {
	err := s.initListener()
	if err != nil {
		return gone.ToError(err)
	}

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.traceInterceptor,
			s.recoveryInterceptor,
		),
	)
	return nil
}

func (s *server) Provide() (*grpc.Server, error) {
	return s.grpcServer, nil
}

func (s *server) register() {
	for _, grpcService := range s.grpcServices {
		s.logger.Infof("Register gRPC service %v", reflect.ValueOf(grpcService).Type().String())
		grpcService.RegisterGrpcServer(s.grpcServer)
	}
}

func (s *server) Start() error {
	s.register()
	if s.tracer == nil {
		go s.server()
	} else {
		s.tracer.Go(s.server)
	}
	return nil
}

func (s *server) server() {
	s.logger.Infof("gRPC server now listen at %s", s.address)
	if err := s.grpcServer.Serve(s.listener); err != nil {
		s.logger.Errorf("failed to serve: %v", err)
	}
}

func (s *server) Stop() error {
	s.grpcServer.Stop()
	return nil
}

func (s *server) traceInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	var traceId string
	traceIdV := metadata.ValueFromIncomingContext(ctx, s.tracerIdKey)
	if len(traceIdV) > 0 {
		traceId = traceIdV[0]
	}

	if s.tracer == nil {
		return handler(ctx, req)
	} else {
		s.tracer.SetTraceId(traceId, func() {
			resp, err = handler(ctx, req)
		})
	}
	return
}

func (s *server) recoveryInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = gone.NewInnerErrorSkip(fmt.Sprintf("panic: %v", e), gone.PanicError, 3)
			s.logger.Errorf("%v", err)
		}
	}()
	return handler(ctx, req)
}
