package balancer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/balancer/strategy"
)

func Load(loader gone.Loader) error {
	var load = gone.OnceLoad(func(loader gone.Loader) error {
		err := loader.Load(&balancer{})
		if err != nil {
			return gone.ToError(err)
		}
		return loader.Load(&strategy.RoundRobinStrategy{})
	})
	return load(loader)
}
