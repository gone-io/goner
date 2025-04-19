package consul

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/spf13/cast"
	"sync"
	"time"
)

var _ g.ServiceRegistry = (*Registry)(nil)
var _ g.ServiceDiscovery = (*Registry)(nil)

type Registry struct {
	gone.Flag

	logger gone.Logger `gone:"*"`
	client *api.Client `gone:"*"`

	servicesMap map[string][]*api.AgentServiceRegistration
	mu          sync.RWMutex // Mutex for thread safety
}

func (r *Registry) Init() (err error) {
	r.servicesMap = make(map[string][]*api.AgentServiceRegistration)
	return
}

const weightKey = "_weight"

// DefaultTTL is the default TTL for service registration
const DefaultTTL = 20 * time.Second

// DefaultHealthCheckInterval is the default interval for health check
const DefaultHealthCheckInterval = 10 * time.Second

func (r *Registry) GetInstances(serviceName string) ([]g.Service, error) {
	list, _, err := r.client.Health().Service(serviceName, "", true, &api.QueryOptions{
		WaitTime: time.Second * 3,
	})
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "get instances failed")
	}

	services := make([]g.Service, 0, len(list))
	for _, s := range list {
		services = append(services,
			g.NewService(
				s.Service.Service,
				s.Service.Address,
				s.Service.Port,
				s.Service.Meta,
				true,
				cast.ToFloat64(s.Service.Meta[weightKey]),
			),
		)
	}
	return services, nil
}

func (r *Registry) Watch(serviceName string) (ch <-chan []g.Service, stop func() error, err error) {
	r.mu.Lock()
	plan, err := watch.Parse(map[string]any{
		"type":    "service",
		"service": serviceName,
	})
	r.mu.Unlock()

	if err != nil {
		return nil, nil, gone.ToErrorWithMsg(err, "parse watch plan failed")
	}

	c := make(chan []g.Service)
	plan.Handler = func(u uint64, data any) {
		instances, err := r.GetInstances(serviceName)
		if err != nil {
			r.logger.Errorf("get instances failed: %v", err)
			return
		}
		c <- instances
	}

	go func() {
		if err := plan.RunWithClientAndHclog(r.client, nil); err != nil {
			r.logger.Errorf("watch plan err: %v", err)
		}
	}()

	return c, func() error {
		r.mu.Lock()
		if plan != nil {
			plan.Stop()
			plan = nil
		}
		r.mu.Unlock()
		close(c)
		return nil
	}, nil
}

func getServiceId(instance g.Service) string {
	return fmt.Sprintf("%s-%s:%d", instance.GetName(), instance.GetIP(), instance.GetPort())
}

func (r *Registry) Register(instance g.Service) error {
	metadata := instance.GetMetadata()
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata[weightKey] = fmt.Sprintf("%f", instance.GetWeight())
	serviceID := getServiceId(instance)

	registration := api.AgentServiceRegistration{
		ID: serviceID, Name: instance.GetName(),
		Address: instance.GetIP(),
		Port:    instance.GetPort(),
		Meta:    metadata,
	}

	checkID := fmt.Sprintf("service:%s", registration.ID)
	registration.Check = &api.AgentServiceCheck{
		CheckID:                        checkID,
		TTL:                            DefaultTTL.String(),
		DeregisterCriticalServiceAfter: "1m",
	}

	if err := r.client.Agent().ServiceRegister(&registration); err != nil {
		return gone.ToErrorWithMsg(err, "register service failed")
	}

	// Start TTL health check
	if err := r.client.Agent().PassTTL(checkID, ""); err != nil {
		// Try to deregister service if health check fails
		_ = r.client.Agent().ServiceDeregister(serviceID)
		return gone.ToErrorWithMsg(err, "failed to pass TTL health check")
	}

	// Start TTL health check goroutine
	go r.ttlHealthCheck(serviceID)
	return nil
}

// ttlHealthCheck maintains the TTL health check for a service
func (r *Registry) ttlHealthCheck(serviceID string) {
	ticker := time.NewTicker(DefaultHealthCheckInterval)
	defer ticker.Stop()

	checkID := fmt.Sprintf("service:%s", serviceID)
	for range ticker.C {
		if err := r.client.Agent().PassTTL(checkID, ""); err != nil {
			return
		}
	}
}

func (r *Registry) Deregister(instance g.Service) error {
	serviceID := getServiceId(instance)
	if err := r.client.Agent().ServiceDeregister(serviceID); err != nil {
		return gone.ToErrorWithMsg(err, "deregister service failed")
	}
	return nil
}
