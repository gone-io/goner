package etcd

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	etcd3 "go.etcd.io/etcd/client/v3"
)

// extractResponseToServices extracts etcd watch response context to service list.
func extractResponseToServices(res *etcd3.GetResponse) ([]g.Service, error) {
	if res == nil || res.Kvs == nil {
		return nil, nil
	}
	var services []g.Service
	for _, kv := range res.Kvs {
		service, err := g.ParseService(string(kv.Value))
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, "parse service failed")
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *Registry) GetInstances(serviceName string) ([]g.Service, error) {
	res, err := r.client.Get(context.Background(), serviceName, etcd3.WithPrefix())
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf(`etcd get failed with key "%s"`, serviceName))
	}
	services, err := extractResponseToServices(res)
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf(`etcd get failed with key "%s"`, serviceName))
	}
	return services, nil
}

func (r *Registry) Watch(serviceName string) (<-chan []g.Service, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.dialTimeout)
	defer cancel()

	// Test connection first.
	if _, err := r.client.Get(ctx, "ping"); err != nil {
		return nil, nil, gone.ToErrorWithMsg(err, "failed to connect to etcd")
	}
	return r.watch(serviceName, etcd3.NewWatcher(r.client))
}

func (r *Registry) watch(serviceName string, watcher etcd3.Watcher) (<-chan []g.Service, func() error, error) {
	wCtx, wCancel := context.WithCancel(context.Background())
	watchCh := watcher.Watch(wCtx, serviceName, etcd3.WithPrefix(), etcd3.WithRev(0))
	if err := watcher.RequestProgress(context.Background()); err != nil {
		wCancel()
		return nil, nil, gone.ToErrorWithMsg(err, "failed to request progress")
	}
	ch := make(chan []g.Service)

	go func() {
		for {
			<-watchCh
			instances, err := r.GetInstances(serviceName)
			if err != nil {
				r.logger.Errorf("get instances failed: %v", err)
				continue
			}
			ch <- instances
		}
	}()

	return ch, func() error {
		wCancel()
		return gone.ToError(watcher.Close())
	}, nil
}
