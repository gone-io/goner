package grpc

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func createListener(s *server) (err error) {
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	return
}

type server struct {
	gone.Flag
	logger       gone.Logger         `gone:"*"`
	grpcServices []Service           `gone:"*"`
	grpcOptions  []grpc.ServerOption `gone:"*"`
	cMuxServer   g.Cmux              `gone:"*" option:"allowNil"`
	tracer       g.Tracer            `gone:"*" option:"allowNil"`
	registry     g.ServiceRegistry   `gone:"*" option:"allowNil"`

	port             int    `gone:"config,server.grpc.port,default=9090"`
	host             string `gone:"config,server.grpc.host,default=0.0.0.0"`
	serviceName      string `gone:"config,server.grpc.service-name"`
	serviceUseSubNet string `gone:"config,server.grpc.service-use-subnet,default=0.0.0.0/0"`
	tracerIdKey      string `gone:"config,server.grpc.x-trace-id-key=X-Trace-Id"`

	grpcServer     *grpc.Server
	listener       net.Listener
	createListener func(*server) error
	unRegService   func() error
}

func (s *server) GonerName() string {
	return "gone-grpc-server"
}

func (s *server) getAddress() string {
	if s.cMuxServer != nil {
		return s.cMuxServer.GetAddress()
	}
	return s.listener.Addr().String()
}

func (s *server) initListener() error {
	if s.cMuxServer != nil {
		s.listener = s.cMuxServer.MatchFor(g.GRPC)
		return nil
	}
	return s.createListener(s)
}
func (s *server) Init() error {
	err := s.initListener()
	if err != nil {
		return gone.ToError(err)
	}

	options := append(s.grpcOptions, grpc.ChainUnaryInterceptor(
		s.traceInterceptor,
		s.recoveryInterceptor,
	))

	s.grpcServer = grpc.NewServer(options...)
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

func (s *server) GetPort() int {
	if s.listener == nil {
		return s.port
	}
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *server) regService() func() error {
	if s.registry != nil {
		if s.serviceName == "" {
			panic("serviceName is empty, please config serviceName by setting key `server.grpc.service-name` value")
		}

		s.logger.Infof("Register gRPC service %v", reflect.ValueOf(s).Type().String())
		ips := g.GetLocalIps()
		port := s.GetPort()

		_, ipnet, err := net.ParseCIDR(s.serviceUseSubNet)
		if err != nil {
			panic(fmt.Sprintf("serviceUseSubNet is invalid, please config serviceUseSubNet by setting key `server.grpc.service-use-subnet` value"))
		}

		for _, ip := range ips {
			if ipnet.Contains(ip) {
				service := g.NewService(s.serviceName, ip.String(), port, nil, true, 100)
				err := s.registry.Register(service)
				if err != nil {
					s.logger.Errorf("register gRPC service %s failed: %v", s.serviceName, err)
					panic(err)
				}
				s.logger.Debugf("Register gRPC service %s success with %s:%d", service.GetName(), service.GetIP(), service.GetPort())
				return func() error {
					return gone.ToError(s.registry.Deregister(service))
				}
			}
		}
		panic(fmt.Sprintf("serviceUseSubNet is invalid, please config serviceUseSubNet by setting key `server.grpc.service-use-subnet` value"))
	}
	return nil
}

func (s *server) Start() error {
	s.register()
	if s.tracer == nil {
		go s.server()
	} else {
		s.tracer.Go(s.server)
	}
	s.unRegService = s.regService()
	return nil
}

func (s *server) server() {
	s.logger.Infof("gRPC server now listen at %s", s.getAddress())
	if err := s.grpcServer.Serve(s.listener); err != nil {
		s.logger.Errorf("failed to serve: %v", err)
	}
}

func (s *server) Stop() error {
	if s.unRegService != nil {
		err := s.unRegService()
		if err != nil {
			s.logger.Errorf("unregister gRPC service %s failed: %v", s.serviceName, err)
		}
	}
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
		return
	}
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
