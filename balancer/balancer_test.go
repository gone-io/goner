package balancer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	mock "github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	gMock "github.com/gone-io/goner/g/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 创建测试用的服务实例
func createTestService(name, ip string, port int) g.Service {
	return g.NewService(name, ip, port, g.Metadata{}, true, 1.0)
}

// TestBalancer_GetInstance 测试GetInstance方法
func TestBalancer_GetInstance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	mockDiscovery := gMock.NewMockServiceDiscovery(ctrl)
	mockStrategy := gMock.NewMockLoadBalanceStrategy(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)
	mockService := gMock.NewMockService(ctrl)

	// 创建balancer实例
	b := &balancer{
		strategy:  mockStrategy,
		discovery: mockDiscovery,
		logger:    mockLogger,
		m:         sync.Map{},
	}

	// 测试场景1: 缓存中没有服务实例，需要从discovery获取
	t.Run("Get instance from discovery when cache is empty", func(t *testing.T) {
		serviceName := "test-service"
		ctx := context.Background()
		instances := []g.Service{mockService}

		// 设置模拟行为
		mockDiscovery.EXPECT().GetInstances(serviceName).Return(instances, nil)
		mockDiscovery.EXPECT().Watch(serviceName).Return(make(<-chan []g.Service), func() error { return nil }, nil)
		mockStrategy.EXPECT().Select(ctx, instances).Return(mockService, nil)
		mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()

		// 执行测试
		service, err := b.GetInstance(ctx, serviceName)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, mockService, service)
	})

	// 测试场景2: 缓存中已有服务实例
	t.Run("Get instance from cache", func(t *testing.T) {
		serviceName := "cached-service"
		ctx := context.Background()
		instances := []g.Service{mockService}

		// 预先设置缓存
		b.m.Store(serviceName, instances)

		// 设置模拟行为
		mockStrategy.EXPECT().Select(ctx, instances).Return(mockService, nil)

		// 执行测试
		service, err := b.GetInstance(ctx, serviceName)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, mockService, service)
	})

	// 测试场景3: 获取服务实例失败
	t.Run("Fail to get instance from discovery", func(t *testing.T) {
		serviceName := "error-service"
		ctx := context.Background()
		expectedErr := errors.New("discovery error")

		// 设置模拟行为
		mockDiscovery.EXPECT().GetInstances(serviceName).Return(nil, expectedErr)

		// 执行测试
		service, err := b.GetInstance(ctx, serviceName)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, service)
	})

	// 测试场景4: 策略选择失败
	t.Run("Strategy select fails", func(t *testing.T) {
		serviceName := "strategy-error-service"
		ctx := context.Background()
		instances := []g.Service{mockService}
		expectedErr := errors.New("strategy error")

		// 预先设置缓存
		b.m.Store(serviceName, instances)

		// 设置模拟行为
		mockStrategy.EXPECT().Select(ctx, instances).Return(nil, expectedErr)

		// 执行测试
		service, err := b.GetInstance(ctx, serviceName)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, service)
	})
}

// TestBalancer_GetInstancesWithCacheAndWatch 测试GetInstancesWithCacheAndWatch方法
func TestBalancer_GetInstancesWithCacheAndWatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟对象
	mockDiscovery := gMock.NewMockServiceDiscovery(ctrl)
	mockStrategy := gMock.NewMockLoadBalanceStrategy(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	// 创建balancer实例
	b := &balancer{
		strategy:  mockStrategy,
		discovery: mockDiscovery,
		logger:    mockLogger,
		m:         sync.Map{},
	}

	// 测试场景1: 缓存中没有服务实例，需要从discovery获取并启动监听
	t.Run("Get instances from discovery and start watching", func(t *testing.T) {
		serviceName := "test-service"
		service1 := createTestService(serviceName, "127.0.0.1", 8080)
		service2 := createTestService(serviceName, "127.0.0.1", 8081)
		instances := []g.Service{service1, service2}

		// 创建用于监听的通道
		ch := make(chan []g.Service)
		stopFunc := func() error { return nil }

		// 设置模拟行为
		mockDiscovery.EXPECT().GetInstances(serviceName).Return(instances, nil)
		mockDiscovery.EXPECT().Watch(serviceName).Return(ch, stopFunc, nil)
		mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		// 执行测试
		result, err := b.GetInstancesWithCacheAndWatch(serviceName)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, instances, result)

		// 模拟服务实例更新
		updatedService := createTestService(serviceName, "127.0.0.1", 8082)
		updatedInstances := []g.Service{updatedService}

		// 发送更新到通道
		go func() {
			ch <- updatedInstances
		}()

		// 等待更新生效
		time.Sleep(100 * time.Millisecond)

		// 验证缓存已更新
		cachedValue, ok := b.m.Load(serviceName)
		assert.True(t, ok)
		cachedInstances, ok := cachedValue.([]g.Service)
		assert.True(t, ok)
		assert.Equal(t, updatedInstances, cachedInstances)
	})

	// 测试场景2: 缓存中已有服务实例
	t.Run("Get instances from cache", func(t *testing.T) {
		serviceName := "cached-service"
		service1 := createTestService(serviceName, "127.0.0.1", 8080)
		service2 := createTestService(serviceName, "127.0.0.1", 8081)
		instances := []g.Service{service1, service2}

		// 预先设置缓存
		b.m.Store(serviceName, instances)

		// 执行测试
		result, err := b.GetInstancesWithCacheAndWatch(serviceName)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, instances, result)
	})

	// 测试场景3: 获取服务实例失败
	t.Run("Fail to get instances from discovery", func(t *testing.T) {
		serviceName := "error-service"
		expectedErr := errors.New("discovery error")

		// 设置模拟行为
		mockDiscovery.EXPECT().GetInstances(serviceName).Return(nil, expectedErr)

		// 执行测试
		result, err := b.GetInstancesWithCacheAndWatch(serviceName)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// 测试场景4: 监听失败
	t.Run("Watch fails", func(t *testing.T) {
		serviceName := "watch-error-service"
		service1 := createTestService(serviceName, "127.0.0.1", 8080)
		instances := []g.Service{service1}
		expectedErr := errors.New("watch error")

		// 设置模拟行为
		mockDiscovery.EXPECT().GetInstances(serviceName).Return(instances, nil)
		mockDiscovery.EXPECT().Watch(serviceName).Return(nil, nil, expectedErr)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		// 执行测试
		result, err := b.GetInstancesWithCacheAndWatch(serviceName)

		// 验证结果 - 即使监听失败，也应该返回实例
		assert.NoError(t, err)
		assert.Equal(t, instances, result)

		time.Sleep(10 * time.Millisecond)
	})
}
