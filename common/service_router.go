package common

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/interfaces/svc"
	"net/http"
)

type serviceRouter struct {
	gone.Flag
	discovery svc.Discovery `gone:"*"`
}

func (s *serviceRouter) GetServiceAddress(serviceName http.Request) (serviceAddress string, err error) {

}
