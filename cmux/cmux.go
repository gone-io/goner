package cmux

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/tracer"
	"github.com/soheilhy/cmux"
	"net"
	"net/http"
	"sync"
	"time"
)

const Name = "cmux"

var load = gone.OnceLoad(func(loader gone.Loader) error {
	err := tracer.Load(loader)
	if err != nil {
		return gone.ToError(err)
	}
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
	once   sync.Once
	cMux   cmux.CMux
	logger gone.Logger   `gone:"*"`
	tracer tracer.Tracer `gone:"*" option:"allowNil"`

	stopFlag bool
	lock     sync.Mutex

	network string `gone:"config,server.network,default=tcp"`
	address string `gone:"config,server.address"`
	host    string `gone:"config,server.host"`
	port    int    `gone:"config,server.port,default=8080"`

	listen func(network, address string) (net.Listener, error)
}

func (s *server) GonerName() string {
	return Name
}

func (s *server) Init() error {
	var err error
	if s.cMux == nil {
		s.once.Do(func() {
			if s.address == "" {
				s.address = fmt.Sprintf("%s:%d", s.host, s.port)
			}
			var listener net.Listener
			listener, err = s.listen(s.network, s.address)
			if err != nil {
				return
			}
			s.cMux = cmux.New(listener)
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
		return s.MatchWithWriters(
			cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
		)
	case g.HTTP1:
		return s.Match(cmux.HTTP1Fast(http.MethodPatch))
	default:
		panic(gone.ToError(fmt.Errorf("unsupport protocol type:%d", protocol)))
	}
}

func (s *server) MatchWithWriters(matcher ...cmux.MatchWriter) net.Listener {
	return s.cMux.MatchWithWriters(matcher...)
}

func (s *server) GetAddress() string {
	return s.address
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

	if s.tracer != nil {
		s.tracer.Go(fn)
	} else {
		go fn()
	}
	<-time.After(20 * time.Millisecond)
	return err
}

func (s *server) Stop() error {
	s.logger.Warnf("cMux server stopping!!")
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stopFlag = true
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
