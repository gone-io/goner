package balancer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/balancer/strategy"
	"github.com/gone-io/goner/g"
)

var load = g.BuildOnceLoadFunc(
	g.L(&balancer{}),
	g.L(&strategy.RoundRobinStrategy{}),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
