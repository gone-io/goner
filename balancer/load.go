package balancer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/balancer/strategy"
	"github.com/gone-io/goner/g"
)

func Load(loader gone.Loader) error {
	return g.BuildLoadFunc(loader, g.L(&balancer{}), g.L(&strategy.RoundRobinStrategy{}))
}
