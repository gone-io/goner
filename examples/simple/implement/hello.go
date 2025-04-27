package implement

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"simple/service"
)

var _ service.Service = (*serviceImpl)(nil)

type serviceImpl struct {
	gone.Flag
}

func (s serviceImpl) SayHello(name string) string {
	return fmt.Sprintf("hello %s", name)
}
