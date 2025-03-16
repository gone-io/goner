package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/cmux"
	"github.com/gone-io/goner/tracer"
	Cmux "github.com/soheilhy/cmux"
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
	s.listener, err = net.Listen("tcp", s.address)
	return
}

type server struct {
	gone.Flag
	httpServer  *http.Server
	logger      gone.Logger     `gone:"gone-logger"`
	httpHandler http.Handler    `gone:"gone-gin-router"`
	cMuxServer  cmux.CMuxServer `gone:"*" option:"allowNil"`
	tracer      tracer.Tracer   `gone:"*" option:"allowNil"`

	controllers []Controller `gone:"*"`

	address  string
	stopFlag bool
	lock     sync.Mutex

	listener          net.Listener
	port              int           `gone:"config,server.port=8080"`
	host              string        `gone:"config,server.host,default=0.0.0.0"`
	maxWaitBeforeStop time.Duration `gone:"config,server.max-wait-before-stop=5s"`

	createListener func(*server) error
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

	s.logger.Infof("Server Listen At http://%s", s.address)
	if s.tracer == nil {
		go s.serve()
	} else {
		s.tracer.Go(s.serve)
	}
	return nil
}

func (s *server) initListener() error {
	if s.cMuxServer != nil {
		s.listener = s.cMuxServer.Match(Cmux.HTTP1Fast(http.MethodPatch))
		s.address = s.cMuxServer.GetAddress()
		return nil
	}

	s.address = fmt.Sprintf("%s:%d", s.host, s.port)
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
