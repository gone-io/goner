package hello

import (
	"examples/simple/internal/interface/service"
	"fmt"
	"github.com/gone-io/gone/v2"
)

var _ service.IService = (*serviceImpl)(nil)

type serviceImpl struct {
	gone.Flag
}

func (s serviceImpl) SayHello(name string) string {
	return fmt.Sprintf("hello %s", name)
}
