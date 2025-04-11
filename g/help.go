package g

import (
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
		panic(gone.ToErrorWithMsg(err, "cannot get ip address"))
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
	f       gone.LoadFunc
}

func L(g gone.Goner, options ...gone.Option) *LoadOp {
	return &LoadOp{
		goner:   g,
		options: options,
	}
}

func F(loadFunc gone.LoadFunc) *LoadOp {
	return &LoadOp{
		f: loadFunc,
	}
}

func BuildOnceLoadFunc(ops ...*LoadOp) gone.LoadFunc {
	return gone.OnceLoad(func(loader gone.Loader) error {
		for _, op := range ops {
			if op.goner != nil {
				err := loader.Load(
					op.goner,
					op.options...,
				)
				if err != nil {
					return gone.ToError(err)
				}
			}
			if op.f != nil {
				err := op.f(loader)
				if err != nil {
					return gone.ToError(err)
				}
			}
		}
		return nil
	})
}
