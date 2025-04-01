package strategy

import (
	"context"
	"math/rand"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var _ g.LoadBalanceStrategy = (*WeightStrategy)(nil)

type WeightStrategy struct {
	gone.Flag
}

// Select 根据权重选择实例
func (w WeightStrategy) Select(ctx context.Context, instances []g.Service) (g.Service, error) {
	if len(instances) == 0 {
		return nil, gone.ToError("no available service instances")
	}

	var totalWeight int64 = 0
	for _, instance := range instances {
		totalWeight += int64(instance.GetWeight() * 100)
	}

	if totalWeight <= 0 {
		return nil, gone.ToError("total weight must be greater than zero")
	}

	randWeight := rand.Int63n(totalWeight)
	var currentWeight int64 = 0
	for _, instance := range instances {
		currentWeight += int64(instance.GetWeight() * 100)
		if randWeight < currentWeight {
			return instance, nil
		}
	}
	return nil, gone.ToError("failed to select instance by weight")
}
