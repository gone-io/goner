package strategy

import (
	"context"
	"sync/atomic"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var _ g.LoadBalanceStrategy = (*RoundRobinStrategy)(nil)

type RoundRobinStrategy struct {
	gone.Flag
	counter uint64
}

func (r *RoundRobinStrategy) Select(ctx context.Context, instances []g.Service) (g.Service, error) {
	if len(instances) == 0 {
		return nil, gone.ToError("no available service instances")
	}

	// 使用原子操作增加计数器并获取当前值
	index := atomic.AddUint64(&r.counter, 1) % uint64(len(instances))

	// 返回选中的实例
	return instances[index], nil
}
