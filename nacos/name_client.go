package nacos

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/interfaces/svc"
	"github.com/gone-io/goner/internal/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Registry struct {
	gone.Flag

	clientConfig  *constant.ClientConfig  `gone:"config,nacos.client"`
	serverConfigs []constant.ServerConfig `gone:"config,nacos.server"`
	serviceName   string                  `gone:"config,nacos.service.name"`
	groupName     string                  `gone:"config,nacos.service.group"`
	clusterName   string                  `gone:"config,nacos.service.clusterName"`

	iClient naming_client.INamingClient
}

func (reg *Registry) Init() (err error) {
	reg.iClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  reg.clientConfig,
			ServerConfigs: reg.serverConfigs,
		},
	)
	if err != nil {
		return gone.ToError(err)
	}
	return gone.ToError(err)
}

func covertString(value any) string {
	if jsonContent, err := json.Marshal(value); err != nil {
		return fmt.Sprint(value)
	} else {
		return string(jsonContent)
	}
}

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (reg *Registry) Register(service svc.Service) (registered svc.Service, err error) {
	serviceName := reg.serviceName
	if service.GetName() != "" {
		serviceName = service.GetName()
	}

	metadata := map[string]string{}
	endpoints := service.GetEndpoints()
	p := vo.BatchRegisterInstanceParam{
		ServiceName: serviceName,
		GroupName:   reg.groupName,
		Instances:   make([]vo.RegisterInstanceParam, 0, len(endpoints)),
	}

	for k, v := range service.GetMetadata() {
		metadata[k] = covertString(v)
	}

	for _, endpoint := range endpoints {
		p.Instances = append(p.Instances, vo.RegisterInstanceParam{
			Ip:          endpoint.Host(),
			Port:        uint64(endpoint.Port()),
			ServiceName: serviceName,
			Metadata:    metadata,
			Weight:      100,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			ClusterName: reg.clusterName,
			GroupName:   reg.groupName,
		})
	}

	if _, err = reg.iClient.BatchRegisterInstance(p); err != nil {
		return
	}

	registered = service

	return
}

// Deregister off-lines and removes `service` from the Registry.
func (reg *Registry) Deregister(service svc.Service) (err error) {
	serviceName := reg.serviceName
	if service.GetName() != "" {
		serviceName = service.GetName()
	}

	for _, endpoint := range service.GetEndpoints() {
		if _, err = reg.iClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          endpoint.Host(),
			Port:        uint64(endpoint.Port()),
			ServiceName: serviceName,
			Ephemeral:   true,
			Cluster:     reg.clusterName,
			GroupName:   reg.groupName,
		}); err != nil {
			return
		}
	}
	return
}

// Search searches and returns services with specified condition.
func (reg *Registry) Search(in svc.SearchInput) (result []svc.Service, err error) {
	if in.Name == "" {
		return nil, gone.ToError("in.Name cannot be empty")
	}
	instances, err := reg.iClient.SelectInstances(vo.SelectInstancesParam{
		GroupName:   reg.groupName,
		Clusters:    []string{reg.clusterName},
		ServiceName: in.Name,
		HealthyOnly: true,
	})

	if err != nil {
		return
	}

	insts := make([]model.Instance, 0, len(instances))

instLoop:
	for _, inst := range instances {
		if len(in.Metadata) > 0 {
			for k, v := range in.Metadata {
				if inst.Metadata[k] != v {
					continue instLoop
				}
			}
		}
		insts = append(insts, inst)
	}

	result = NewServicesFromInstances(insts)
	return
}

type endpoint struct {
	host string
	port int
}

// Host returns the IPv4/IPv6 address of a service.
func (e endpoint) Host() string {
	return e.host
}

// Port returns the port of a service.
func (e endpoint) Port() int {
	return e.port
}

// String formats and returns the Endpoint as a string.
func (e endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.host, e.port)
}

type service struct {
	name      string
	endpoints svc.Endpoints
	Metadata  svc.Metadata
}

func (s service) GetName() string {
	return s.name
}

func (s service) GetMetadata() svc.Metadata {
	return s.Metadata
}

func (s service) GetEndpoints() svc.Endpoints {
	return s.endpoints
}

func NewServicesFromInstances(insts []model.Instance) []svc.Service {
	services := make([]svc.Service, 0, len(insts))
	for _, inst := range insts {
		metadata := make(svc.Metadata)
		for k, v := range inst.Metadata {
			metadata[k] = v
		}
		services = append(services, service{
			name: inst.ServiceName,
			endpoints: []svc.Endpoint{
				endpoint{
					host: inst.Ip,
					port: int(inst.Port),
				},
			},
			Metadata: metadata,
		})
	}
	return services
}

// Watch watches specified condition changes.
func (reg *Registry) Watch(serviceName string) (svc.Watcher, error) {
	c := reg.iClient

	w := newWatcher()
	param := &vo.SubscribeParam{
		ServiceName:       serviceName,
		GroupName:         reg.groupName,
		Clusters:          []string{reg.clusterName},
		SubscribeCallback: w.Push,
	}

	w.SetCloseFunc(func() error {
		return c.Unsubscribe(param)
	})

	if err := c.Subscribe(param); err != nil {
		return nil, gone.ToError(err)
	}
	return w, nil
}

// watchEvent
type watchEvent struct {
	Services []model.Instance
	Err      error
}

func newWatcher() *watcher {
	return &watcher{
		event: make(chan *watchEvent, 10),
	}
}

type watcher struct {
	close func() error
	event chan *watchEvent
}

func (w *watcher) Push(services []model.Instance, err error) {
	w.event <- &watchEvent{
		Services: services,
		Err:      err,
	}
}

func (w *watcher) SetCloseFunc(close func() error) {
	w.close = close
}

func (w *watcher) Proceed() (services []svc.Service, err error) {
	e, ok := <-w.event
	if !ok || e == nil {
		err = gone.ToError(err)
		return
	}
	if e.Err != nil {
		err = gone.ToError(e.Err)
		return
	}
	services = NewServicesFromInstances(e.Services)
	return
}

// Close closes the watcher.
func (w *watcher) Close() error {
	if w.close != nil {
		return w.close()
	}
	return nil
}
