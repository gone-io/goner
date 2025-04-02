package balancer

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"sync"
)

var _ g.LoadBalancer = (*balancer)(nil)

type balancer struct {
	gone.Flag
	strategy  g.LoadBalanceStrategy `gone:"*"`
	discovery g.ServiceDiscovery    `gone:"*"`
	logger    gone.Logger           `gone:"*"`
	m         sync.Map
}

func (b *balancer) GetInstance(ctx context.Context, serviceName string) (g.Service, error) {
	var instances []g.Service
	var err error
	value, ok := b.m.Load(serviceName)
	if !ok {
		instances, err = b.GetInstancesWithCacheAndWatch(serviceName)
		if err != nil {
			return nil, gone.ToError(err)
		}
	} else {
		instances, ok = value.([]g.Service)
	}

	return b.strategy.Select(ctx, instances)
}

func (b *balancer) GetInstancesWithCacheAndWatch(serviceName string) ([]g.Service, error) {
	if value, ok := b.m.Load(serviceName); ok {
		return value.([]g.Service), nil
	}

	instances, err := b.discovery.GetInstances(serviceName)
	if err != nil {
		return nil, gone.ToError(err)
	}
	b.m.Store(serviceName, instances)
	go func() {
		defer g.Recover(b.logger)

		ch, stop, err := b.discovery.Watch(serviceName)
		if err != nil {
			b.logger.Errorf("balancer watch %s err: %v", serviceName, err)
			return
		}
		defer stop()
		for {
			select {
			case <-ch:
				instances, err := b.discovery.GetInstances(serviceName)
				if err != nil {
					continue
				}
				b.m.Store(serviceName, instances)
			}
		}
	}()
	return instances, nil
}
