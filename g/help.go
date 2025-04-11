package g

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"net"
)

func Recover(logger gone.Logger) {
	if r := recover(); r != nil {
		logger.Errorf(
			"panic: %v, %s",
			r,
			gone.PanicTrace(2, 1),
		)
	}
}

func GetLocalIps() []net.IP {
	if addrs, err := net.InterfaceAddrs(); err != nil {
		panic(fmt.Sprintf("cannot get ip addresss: %v", err))
	} else {
		var ips []net.IP
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
		return ips
	}
}

type LoadOp struct {
	goner   gone.Goner
	options []gone.Option
}

func L(g gone.Goner, options ...gone.Option) *LoadOp {
	return &LoadOp{
		goner:   g,
		options: options,
	}
}

func BuildLoadFunc(loader gone.Loader, ops ...*LoadOp) error {
	return gone.OnceLoad(func(loader gone.Loader) error {
		for _, op := range ops {
			err := loader.Load(
				op.goner,
				op.options...,
			)
			if err != nil {
				return gone.ToError(err)
			}
		}
		return nil
	})(loader)
}
