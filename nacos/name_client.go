package nacos

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var _ g.ServiceRegistry = (*Registry)(nil)
var _ g.ServiceDiscovery = (*Registry)(nil)

type Registry struct {
	gone.Flag

	logger        gone.Logger             `gone:"*"`
	clientConfig  constant.ClientConfig   `gone:"config,nacos.client"`
	serverConfigs []constant.ServerConfig `gone:"config,nacos.server"`
	groupName     string                  `gone:"config,nacos.service.group"`
	clusterName   string                  `gone:"config,nacos.service.clusterName"`

	iClient naming_client.INamingClient
}

func (reg *Registry) Init() (err error) {
	if reg.iClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &reg.clientConfig,
			ServerConfigs: reg.serverConfigs,
		},
	); err != nil {
		return gone.ToError(err)
	}
	return nil
}

func (reg *Registry) GetInstances(serviceName string) ([]g.Service, error) {
	services, err := reg.iClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		Clusters:    []string{reg.clusterName},
		GroupName:   reg.groupName,
	})

	if err != nil {
		return nil, gone.ToError(err)
	}

	instances := make([]g.Service, 0, len(services))
	for _, service := range services {
		instances = append(instances,
			g.NewService(
				service.ServiceName,
				service.Ip,
				int(service.Port),
				service.Metadata,
				service.Healthy,
				service.Weight,
			),
		)
	}
	return instances, nil
}

func (reg *Registry) Watch(serviceName string) (<-chan []g.Service, func() error, error) {
	ch := make(chan []g.Service)
	SubscribeCallback := func(services []model.Instance, err error) {
		if err != nil {
			reg.logger.Errorf("SubscribeCallback err: %v", err)
			return
		}

		reg.logger.Debugf("SubscribeCallback result: %#+v", services)

		instances := make([]g.Service, 0, len(services))
		for _, service := range services {
			instances = append(instances,
				g.NewService(
					service.ServiceName,
					service.Ip,
					int(service.Port),
					service.Metadata,
					service.Healthy,
					service.Weight,
				),
			)
		}
		ch <- instances
	}

	if err := reg.iClient.Subscribe(&vo.SubscribeParam{
		ServiceName:       serviceName,
		GroupName:         reg.groupName,
		Clusters:          []string{reg.clusterName},
		SubscribeCallback: SubscribeCallback,
	}); err != nil {
		return nil, nil, gone.ToError(err)
	}

	return ch, func() error {
		err := reg.iClient.Unsubscribe(&vo.SubscribeParam{
			ServiceName:       serviceName,
			GroupName:         reg.groupName,
			Clusters:          []string{reg.clusterName},
			SubscribeCallback: SubscribeCallback,
		})
		return gone.ToError(err)
	}, nil
}

func (reg *Registry) Register(instance g.Service) error {
	success, err := reg.iClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          instance.GetIP(),
		Port:        uint64(instance.GetPort()),
		ServiceName: instance.GetName(),
		Weight:      instance.GetWeight(),
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    instance.GetMetadata(),
		ClusterName: reg.clusterName,
		GroupName:   reg.groupName,
	})
	if err != nil {
		return gone.ToError(err)
	}
	if !success {
		return gone.ToError(fmt.Sprintf("Register %#+v failed", instance))
	}
	return nil
}

func (reg *Registry) Deregister(instance g.Service) error {

	success, err := reg.iClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          instance.GetIP(),
		Port:        uint64(instance.GetPort()),
		ServiceName: instance.GetName(),
		Ephemeral:   true,
		Cluster:     reg.clusterName,
		GroupName:   reg.groupName,
	})
	if err != nil {
		return gone.ToError(err)
	}
	if !success {
		return gone.ToError(fmt.Sprintf("Deregister %#+v failed", instance))
	}
	return nil
}
