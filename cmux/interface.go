package cmux

import (
	"github.com/soheilhy/cmux"
	"net"
)

//go:generate mockgen -package=cmux -destination=./net_Listener_mock_test.go net Listener,Conn

// CMuxServer cMux service，Used to multiplex the same port to listen for multiple protocols，ref：https://pkg.go.dev/github.com/soheilhy/cmux
type CMuxServer interface {
	Match(matcher ...cmux.Matcher) net.Listener
	MatchWithWriters(matcher ...cmux.MatchWriter) net.Listener
	GetAddress() string
}

// Server cumx 服务，用于复用同一端口监听多种协议，参考文档：https://pkg.go.dev/github.com/soheilhy/cmux
type Server = CMuxServer
