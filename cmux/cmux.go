package cmux

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/soheilhy/cmux"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const Name = "cmux"

var load = gone.OnceLoad(func(loader gone.Loader) error {
	return loader.Load(
		&server{listen: net.Listen},
		gone.IsDefault(new(CMuxServer)),
		gone.HighStartPriority(),
	)
})

func Load(loader gone.Loader) error {
	return load(loader)
}

type server struct {
	gone.Flag
	once     sync.Once
	cMux     cmux.CMux
	logger   gone.Logger       `gone:"*"`
	tracer   g.Tracer          `gone:"*" option:"allowNil"`
	registry g.ServiceRegistry `gone:"*" option:"allowNil"`

	network          string `gone:"config,server.network,default=tcp"`
	address          string `gone:"config,server.address"`
	host             string `gone:"config,server.host"`
	port             int    `gone:"config,server.port,default=8080"`
	serviceName      string `gone:"config,server.service-name"`
	serviceUseSubNet string `gone:"config,server.service-use-subnet,default=0.0.0.0/0"`

	stopFlag bool
	lock     sync.Mutex
	listener net.Listener

	listen       func(network, address string) (net.Listener, error)
	unRegService func() error

	metadata g.Metadata
}

func (s *server) GonerName() string {
	return Name
}

func (s *server) Init() error {
	s.metadata = make(g.Metadata)

	var err error
	if s.cMux == nil {
		s.once.Do(func() {
			if s.address == "" {
				s.address = fmt.Sprintf("%s:%d", s.host, s.port)
			}
			s.listener, err = s.listen(s.network, s.address)
			if err != nil {
				return
			}
			s.cMux = cmux.New(s.listener)
		})
	}
	return err
}

func (s *server) Match(matcher ...cmux.Matcher) net.Listener {
	return s.cMux.Match(matcher...)
}

func (s *server) MatchFor(protocol g.ProtocolType) net.Listener {
	switch protocol {
	case g.GRPC:
		s.metadata["grpc"] = "true"
		return s.MatchWithWriters(
			cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
		)
	case g.HTTP1:
		s.metadata["http1"] = "true"
		return s.Match(cmux.HTTP1Fast(http.MethodPatch))
	default:
		panic(gone.ToError(fmt.Errorf("unsupport protocol type:%d", protocol)))
	}
}

func (s *server) MatchWithWriters(matcher ...cmux.MatchWriter) net.Listener {
	return s.cMux.MatchWithWriters(matcher...)
}

func (s *server) regService() func() error {
	if s.registry != nil {
		if s.serviceName == "" {
			panic("serviceName is empty, please config serviceName by setting key `server.grpc.service-name` value")
		}

		s.logger.Infof("Register gRPC service %v", reflect.ValueOf(s).Type().String())
		ips := g.GetLocalIps()
		port := s.getPort()

		_, ipnet, err := net.ParseCIDR(s.serviceUseSubNet)
		if err != nil {
			panic(fmt.Sprintf("serviceUseSubNet is invalid, please config serviceUseSubNet by setting key `server.grpc.service-use-subnet` value"))
		}

		for _, ip := range ips {
			if ipnet.Contains(ip) {
				service := g.NewService(s.serviceName, ip.String(), port, s.metadata, true, 100)
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

func (s *server) GetAddress() string {
	return s.listener.Addr().String()
}

func (s *server) getPort() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *server) Start() error {
	s.stopFlag = false
	var err error
	var mutex sync.Mutex

	fn := func() {
		mutex.Lock()
		defer mutex.Unlock()
		err = s.cMux.Serve()
		s.processStartError(err)
	}

	s.logger.Infof("cMux server(%#v) listen on: %s", s.metadata, s.GetAddress())
	if s.tracer != nil {
		s.tracer.Go(fn)
	} else {
		go fn()
	}
	s.unRegService = s.regService()
	<-time.After(20 * time.Millisecond)
	return err
}

func (s *server) Stop() error {
	s.logger.Warnf("cMux server stopping!!")
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stopFlag = true
	if s.unRegService != nil {
		err := s.unRegService()
		if err != nil {
			s.logger.Errorf("unregister service error: %v", err)
		}
	}
	s.cMux.Close()
	return nil
}

func (s *server) processStartError(err error) {
	if err != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
		if s.stopFlag {
			s.logger.Errorf("cMux Serve() err:%v", err)
		} else {
			s.logger.Warnf("cMux Serve() err:%v", err)
		}
	}
}
