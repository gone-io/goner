package grpc

import (
	"fmt"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

var _ resolver.Builder = (*resolverBuilder)(nil)

func NewResolverBuilder(discovery g.ServiceDiscovery, logger gone.Logger) resolver.Builder {
	return &resolverBuilder{discovery: discovery, logger: logger}
}

type resolverBuilder struct {
	discovery g.ServiceDiscovery
	logger    gone.Logger
}

func (b *resolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &discoveryResolver{
		discovery:   b.discovery,
		logger:      b.logger,
		cc:          cc,
		serviceName: target.Endpoint(),
	}

	ch, stop, err := b.discovery.Watch(target.Endpoint())
	if err != nil {
		return nil, err
	}

	r.stop = stop

	r.updateCh = ch

	// Initial resolution
	instances, err := b.discovery.GetInstances(target.Endpoint())
	if err != nil {
		return nil, err
	}
	r.updateState(instances)

	// Start watching for updates
	go r.watch()

	return r, nil
}

func (b *resolverBuilder) Scheme() string {
	return "dns"
}

type discoveryResolver struct {
	discovery   g.ServiceDiscovery
	cc          resolver.ClientConn
	serviceName string
	stop        func() error
	updateCh    <-chan []g.Service
	logger      gone.Logger
}

func (r *discoveryResolver) ResolveNow(resolver.ResolveNowOptions) {
	instances, err := r.discovery.GetInstances(r.serviceName)
	if err != nil {
		r.logger.Errorf("discoveryResolver ResolveNow get instances err: %v", err)
		return
	}
	r.updateState(instances)
}

func (r *discoveryResolver) Close() {
	if r.stop != nil {
		if err := r.stop(); err != nil {
			r.logger.Errorf("discoveryResolver close err: %v", err)
		}
	}
}

func (r *discoveryResolver) watch() {
	for services := range r.updateCh {
		r.updateState(services)
	}
}

func (r *discoveryResolver) updateState(services []g.Service) {
	addresses := make([]resolver.Address, 0, len(services))
	for _, svc := range services {
		addresses = append(addresses, resolver.Address{
			Addr:       fmt.Sprintf("%s:%d", svc.GetIP(), svc.GetPort()),
			ServerName: svc.GetName(),
			Attributes: attributes.New("weight", svc.GetWeight()),
		})
	}

	err := r.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
	if err != nil {
		r.logger.Errorf("discoveryResolver updateState err: %v", err)
	}
}
