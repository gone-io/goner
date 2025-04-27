package balancer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/balancer/strategy"
	"github.com/gone-io/goner/g"
)

const Strategy = "strategy"

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&balancer{}).
		MustLoad(&strategy.RoundRobinStrategy{}, gone.Name(Strategy))
	return nil
}

func LoadRandomStrategy(loader gone.Loader) error {
	return loader.Load(&strategy.RandomStrategy{}, gone.Name(Strategy), gone.ForceReplace())
}

func LoadWeightStrategy(loader gone.Loader) error {
	return loader.Load(&strategy.WeightStrategy{}, gone.Name(Strategy), gone.ForceReplace())
}

func LoadCustomerStrategy[T interface {
	g.LoadBalanceStrategy
	gone.Goner
}](strategy T) (gone.Goner, gone.Option, gone.Option) {
	return strategy, gone.Name(Strategy), gone.ForceReplace()
}
