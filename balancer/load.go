package balancer

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/balancer/strategy"
	"github.com/gone-io/goner/g"
)

const Strategy = "strategy"

var load = g.BuildOnceLoadFunc(
	g.L(&balancer{}),
	g.L(&strategy.RoundRobinStrategy{}, gone.Name(Strategy)),
)

func Load(loader gone.Loader) error {
	return load(loader)
}

var loadRandomStrategy = g.BuildOnceLoadFunc(
	g.L(&strategy.RandomStrategy{}, gone.Name(Strategy), gone.ForceReplace()),
)

func LoadRandomStrategy(loader gone.Loader) error {
	return loadRandomStrategy(loader)
}

var loadWeightStrategy = g.BuildOnceLoadFunc(
	g.L(&strategy.WeightStrategy{}, gone.Name(Strategy), gone.ForceReplace()),
)

func LoadWeightStrategy(loader gone.Loader) error {
	return loadWeightStrategy(loader)
}

func LoadCustomerStrategy[T interface {
	g.LoadBalanceStrategy
	gone.Goner
}](strategy T) (gone.Goner, gone.Option, gone.Option) {
	return strategy, gone.Name(Strategy), gone.ForceReplace()
}
