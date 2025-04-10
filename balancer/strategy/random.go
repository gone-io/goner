package strategy

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"math/rand/v2"
)

var _ g.LoadBalanceStrategy = (*RandomStrategy)(nil)

type RandomStrategy struct {
	gone.Flag
}

func (r *RandomStrategy) Select(ctx context.Context, instances []g.Service) (g.Service, error) {
	if len(instances) == 0 {
		return nil, gone.ToError("no available service instances")
	}
	n := rand.IntN(len(instances))
	return instances[n], nil
}
