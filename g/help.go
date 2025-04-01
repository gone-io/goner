package g

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"net"
)

func Recover(logger gone.Logger) {
	if r := recover(); r != nil {
		logger.Errorf("panic: %v, %s",
			r,
			gone.PanicTrace(2, 1),
		)
	}
}

func GetLocalIps() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("cannot get ip addresss: %v", err))
	}
	var ips []net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP)
		}
	}
	return ips
}
