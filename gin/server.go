package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"net"
	"net/http"
	"sync"
	"time"
)

func NewGinServer() (gone.Goner, gone.Option) {
	s := server{
		createListener: createListener,
	}
	return &s, gone.MediumStartPriority()
}

func createListener(s *server) (err error) {
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	return
}

func (s *server) getAddress() string {
	if s.cMuxServer != nil {
		return s.cMuxServer.GetAddress()
	}
	return s.listener.Addr().String()
}

func (s *server) getPort() int {
	if s.listener == nil {
		return s.port
	}
	return s.listener.Addr().(*net.TCPAddr).Port
}

type server struct {
	gone.Flag
	httpServer  *http.Server
	logger      gone.Logger       `gone:"gone-logger"`
	httpHandler http.Handler      `gone:"gone-gin-router"`
	cMuxServer  g.Cmux            `gone:"*" option:"allowNil"`
	tracer      g.Tracer          `gone:"*" option:"allowNil"`
	registry    g.ServiceRegistry `gone:"*" option:"allowNil"`

	controllers []Controller `gone:"*"`

	stopFlag bool
	lock     sync.Mutex

	listener          net.Listener
	port              int           `gone:"config,server.port=8080"`
	host              string        `gone:"config,server.host,default=0.0.0.0"`
	serviceName       string        `gone:"config,server.service-name"`
	serviceUseSubNet  string        `gone:"config,server.service-use-subnet,default=0.0.0.0/0"`
	maxWaitBeforeStop time.Duration `gone:"config,server.max-wait-before-stop=5s"`

	createListener func(*server) error
	unRegService   func() error
}

func (s *server) GonerName() string {
	return IdGoneGin
}

func (s *server) Start() error {
	err := s.mount()
	if err != nil {
		return err
	}
	err = s.initListener()
	if err != nil {
		return err
	}

	s.stopFlag = false
	s.httpServer = &http.Server{
		Handler: s.httpHandler,
	}

	s.logger.Infof("Server Listen At %s", s.getAddress())
	if s.tracer == nil {
		go s.serve()
	} else {
		s.tracer.Go(s.serve)
	}
	s.unRegService = s.regService()
	return nil
}

func (s *server) regService() func() error {
	if s.cMuxServer == nil && s.registry != nil {
		if s.serviceName == "" {
			panic("serviceName is empty, please config serviceName by setting key `server.grpc.service-name` value")
		}

		ips := g.GetLocalIps()
		port := s.getPort()

		_, ipnet, err := net.ParseCIDR(s.serviceUseSubNet)
		if err != nil {
			panic(fmt.Sprintf("serviceUseSubNet is invalid, please config serviceUseSubNet by setting key `server.grpc.service-use-subnet` value"))
		}

		for _, ip := range ips {
			if ipnet.Contains(ip) {
				service := g.NewService(s.serviceName, ip.String(), port, g.Metadata{"http1": "true"}, true, 100)
				err := s.registry.Register(service)
				if err != nil {
					s.logger.Errorf("register gRPC service %s failed: %v", s.serviceName, err)
					panic(err)
				}
				s.logger.Debugf("Register http service success with name `%s` at %s:%d", service.GetName(), service.GetIP(), service.GetPort())
				return func() error {
					return gone.ToError(s.registry.Deregister(service))
				}
			}
		}
		panic(fmt.Sprintf("serviceUseSubNet is invalid, please config serviceUseSubNet by setting key `server.grpc.service-use-subnet` value"))
	}
	return nil
}

func (s *server) initListener() error {
	if s.cMuxServer != nil {
		s.listener = s.cMuxServer.MatchFor(g.HTTP1)
		return nil
	}
	return s.createListener(s)
}

func (s *server) serve() {
	if err := s.httpServer.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.processServeError(err)
	}
}

func (s *server) processServeError(err error) {
	s.lock.Lock()
	if !s.stopFlag {
		s.logger.Errorf("http server error: %v", err)
		panic(err)
	} else {
		s.logger.Warnf("http server error: %v", err)
	}
	s.lock.Unlock()
}

func (s *server) Stop() (err error) {
	s.logger.Warnf("gin server stopping!!")
	if nil == s.httpServer {
		return nil
	}

	s.lock.Lock()
	s.stopFlag = true
	s.lock.Unlock()
	if s.unRegService != nil {
		err = s.unRegService()
		if err != nil {
			s.logger.Errorf("unregister service error: %v", err)
		}
	}
	s.stop()
	return
}

func (s *server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxWaitBeforeStop)
	defer cancel()

	// 关闭服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Errorf("Server forced to shutdown: %v\n", err)
	}
}

// 挂载路由
func (s *server) mount() error {
	if len(s.controllers) == 0 {
		s.logger.Warnf("There is no controller working")
	}

	for _, c := range s.controllers {
		err := c.Mount()
		if err != nil {
			s.logger.Errorf("controller mount err:%v", err)
			return err
		}
	}
	return nil
}
