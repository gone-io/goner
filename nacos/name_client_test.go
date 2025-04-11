package nacos

import (
	"errors"
	mock "github.com/gone-io/gone/mock/v2"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
)

// 测试 Register 方法
func TestRegister(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockClient := NewMockINamingClient(controller)
	reg := &Registry{
		iClient: mockClient,
	}

	instance := g.NewService("test-service", "127.0.0.1", 8080, map[string]string{}, true, 1.0)

	mockClient.EXPECT().RegisterInstance(gomock.Any()).Return(true, nil)
	err := reg.Register(instance)
	assert.NoError(t, err)

	// 模拟注册失败
	mockClient.EXPECT().RegisterInstance(gomock.Any()).Return(false, nil)
	err = reg.Register(instance)
	assert.Error(t, err)

	// 模拟注册时发生错误
	mockClient.EXPECT().RegisterInstance(gomock.Any()).Return(false, errors.New("mock error"))
	err = reg.Register(instance)
	assert.Error(t, err)
}

// 测试 Deregister 方法
func TestDeregister(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockClient := NewMockINamingClient(controller)
	reg := &Registry{
		iClient: mockClient,
	}

	instance := g.NewService("test-service", "127.0.0.1", 8080, map[string]string{}, true, 1.0)

	// 模拟成功注销
	mockClient.EXPECT().DeregisterInstance(gomock.Any()).Return(true, nil)
	err := reg.Deregister(instance)
	assert.NoError(t, err)

	// 模拟注销失败
	mockClient.EXPECT().DeregisterInstance(gomock.Any()).Return(false, nil)
	err = reg.Deregister(instance)
	assert.Error(t, err)

	// 模拟注销时发生错误
	mockClient.EXPECT().DeregisterInstance(gomock.Any()).Return(false, errors.New("mock error"))
	err = reg.Deregister(instance)
	assert.Error(t, err)
}

// 测试 GetInstances 方法
func TestGetInstances(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockClient := NewMockINamingClient(controller)
	reg := &Registry{
		iClient: mockClient,
	}

	// 模拟返回实例列表
	mockClient.EXPECT().SelectInstances(gomock.Any()).Return([]model.Instance{
		{
			ServiceName: "test-service",
			Ip:          "127.0.0.1",
			Port:        8080,
			Metadata:    map[string]string{},
			Healthy:     true,
			Weight:      1.0,
		},
	}, nil)

	instances, err := reg.GetInstances("test-service")
	assert.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, "test-service", instances[0].GetName())

	// 模拟获取实例时发生错误
	mockClient.EXPECT().SelectInstances(gomock.Any()).Return(nil, errors.New("mock error"))
	_, err = reg.GetInstances("test-service")
	assert.Error(t, err)
}

// 测试 Watch 方法
func TestWatch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockClient := NewMockINamingClient(controller)
	mockClient.EXPECT().
		Subscribe(gomock.Any()).
		Do(func(param *vo.SubscribeParam) {
			callback := param.SubscribeCallback
			go callback([]model.Instance{
				{
					ServiceName: "test-service",
					Ip:          "127.0.0.1",
					Port:        8080,
					Metadata:    map[string]string{},
					Healthy:     true,
					Weight:      1.0,
				},
			}, nil)
		}).
		Return(nil)

	mockClient.EXPECT().Unsubscribe(gomock.Any())

	logger := mock.NewMockLogger(controller)
	logger.EXPECT().Debugf(gomock.Any(), gomock.Any())

	reg := &Registry{
		iClient: mockClient,
		logger:  logger,
	}

	ch, stop, err := reg.Watch("test-service")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, stop())
	}()

	// 模拟订阅回调

	instances := <-ch
	assert.Len(t, instances, 1)
	assert.Equal(t, "test-service", instances[0].GetName())
}
