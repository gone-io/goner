package g

import (
	"net"
)

type ProtocolType int

const (
	GRPC  ProtocolType = 0x01
	HTTP1 ProtocolType = 0x01 << 1
)

type Cmux interface {
	MatchFor(ProtocolType) net.Listener
	GetAddress() string
}
