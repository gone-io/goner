package implement

import (
	"examples/simple/service"
	"fmt"
	"github.com/gone-io/gone/v2"
)

var _ service.Service = (*serviceImpl)(nil)

type serviceImpl struct {
	gone.Flag
}

func (s serviceImpl) SayHello(name string) string {
	return fmt.Sprintf("hello %s", name)
}
