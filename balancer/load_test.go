package balancer

import (
	"github.com/gone-io/goner/balancer/strategy"
	"github.com/gone-io/goner/g"
	gMock "github.com/gone-io/goner/g/mock"
	"testing"

	mock "github.com/gone-io/gone/mock/v2"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestLoad_Success 测试Load函数成功加载balancer和默认策略
func TestLoad_Success(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)

	// 设置期望：检查是否已加载
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：成功加载balancer
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&balancer{}), gomock.Any()).
		Return(nil)

	// 设置期望：成功加载RoundRobinStrategy
	mockLoader.EXPECT().
		Load(gomock.Any(), gomock.Any()).
		Return(nil)

	// 执行测试
	err := Load(mockLoader)

	// 验证结果
	assert.NoError(t, err)
}

// TestLoad_BalancerError 测试加载balancer失败的场景
func TestLoad_BalancerError(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)

	// 设置期望：检查是否已加载
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：加载balancer失败
	expectedErr := gone.NewError(400, "mock balancer error", 500)
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&balancer{}), gomock.Any()).
		Return(expectedErr)

	// 执行测试
	err := Load(mockLoader)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestLoad_StrategyError 测试加载策略失败的场景
func TestLoad_StrategyError(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)

	// 设置期望：检查是否已加载
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：成功加载balancer
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&balancer{}), gomock.Any()).
		Return(nil)

	// 设置期望：加载策略失败
	expectedErr := gone.NewError(400, "mock strategy error", 500)
	mockLoader.EXPECT().
		Load(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	// 执行测试
	err := Load(mockLoader)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestLoad_AlreadyLoaded 测试已经加载过的场景
func TestLoad_AlreadyLoaded(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)

	// 设置期望：检查是否已加载，返回true表示已加载
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(true)

	// 执行测试
	err := Load(mockLoader)

	// 验证结果
	assert.NoError(t, err)
}

func TestLoadRandomStrategy_Success(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	discovery := gMock.NewMockServiceDiscovery(controller)
	provider := gone.WrapFunctionProvider(func(tagConf string, param struct{}) (g.ServiceDiscovery, error) {
		return discovery, nil
	})

	gone.
		NewApp(Load, LoadRandomStrategy).
		Load(provider).
		Run(func(s g.LoadBalanceStrategy) {
			_, ok := s.(*strategy.RandomStrategy)
			assert.True(t, ok)
		})
}

func TestLoadWeightStrategy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	discovery := gMock.NewMockServiceDiscovery(controller)
	provider := gone.WrapFunctionProvider(func(tagConf string, param struct{}) (g.ServiceDiscovery, error) {
		return discovery, nil
	})

	gone.
		NewApp(Load, LoadWeightStrategy).
		Load(provider).
		Run(func(s g.LoadBalanceStrategy) {
			_, ok := s.(*strategy.WeightStrategy)
			assert.True(t, ok)
		})
}

func TestLoadCustomerStrategy(t *testing.T) {
	LoadCustomerStrategy(&strategy.RandomStrategy{})
}
