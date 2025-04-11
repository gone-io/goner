package strategy

import (
	"context"
	"testing"

	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
)

// 创建测试用的服务实例
func createTestService(name, ip string, port int, weight float64) g.Service {
	return g.NewService(name, ip, port, g.Metadata{}, true, weight)
}

// TestRandomStrategy_Select 测试RandomStrategy的Select方法
func TestRandomStrategy_Select(t *testing.T) {
	// 创建策略实例
	strategy := &RandomStrategy{}

	// 测试场景1: 正常选择服务实例
	t.Run("Select instance successfully", func(t *testing.T) {
		// 创建测试数据
		service1 := createTestService("test-service", "127.0.0.1", 8080, 1.0)
		service2 := createTestService("test-service", "127.0.0.1", 8081, 1.0)
		instances := []g.Service{service1, service2}
		ctx := context.Background()

		// 执行测试
		service, err := strategy.Select(ctx, instances)

		// 验证结果
		assert.NoError(t, err)
		assert.Contains(t, instances, service)
	})

	// 测试场景2: 服务实例列表为空
	t.Run("Empty instances list", func(t *testing.T) {
		// 创建空的实例列表
		instances := []g.Service{}
		ctx := context.Background()

		// 执行测试
		service, err := strategy.Select(ctx, instances)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "no available service instances")
	})

	// 测试场景3: 服务实例列表为nil
	t.Run("Nil instances list", func(t *testing.T) {
		// 创建nil实例列表
		var instances []g.Service
		ctx := context.Background()

		// 执行测试
		service, err := strategy.Select(ctx, instances)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "no available service instances")
	})
}

// TestRoundRobinStrategy_Select 测试RoundRobinStrategy的Select方法
func TestRoundRobinStrategy_Select(t *testing.T) {
	// 创建策略实例
	strategy := &RoundRobinStrategy{}

	// 测试场景1: 正常轮询选择服务实例
	t.Run("Round robin selection", func(t *testing.T) {
		// 创建测试数据
		service1 := createTestService("test-service", "127.0.0.1", 8080, 1.0)
		service2 := createTestService("test-service", "127.0.0.1", 8081, 1.0)
		instances := []g.Service{service1, service2}
		ctx := context.Background()

		// 第一次调用
		service, err := strategy.Select(ctx, instances)
		assert.NoError(t, err)
		firstSelected := service

		// 第二次调用
		service, err = strategy.Select(ctx, instances)
		assert.NoError(t, err)
		secondSelected := service

		// 第三次调用
		service, err = strategy.Select(ctx, instances)
		assert.NoError(t, err)
		thirdSelected := service

		// 验证轮询效果
		assert.Equal(t, firstSelected, thirdSelected, "第三次选择应该回到第一个实例")
		assert.NotEqual(t, firstSelected, secondSelected, "第一次和第二次选择应该不同")
	})

	// 测试场景2: 服务实例列表为空
	t.Run("Empty instances list", func(t *testing.T) {
		// 创建空的实例列表
		instances := []g.Service{}
		ctx := context.Background()

		// 执行测试
		service, err := strategy.Select(ctx, instances)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "no available service instances")
	})
}

// TestWeightStrategy_Select 测试WeightStrategy的Select方法
func TestWeightStrategy_Select(t *testing.T) {
	// 创建策略实例
	strategy := &WeightStrategy{}

	// 测试场景1: 根据权重选择服务实例
	t.Run("Weight based selection", func(t *testing.T) {
		// 创建测试数据，service1的权重是service2的9倍
		service1 := createTestService("test-service", "127.0.0.1", 8080, 9.0)
		service2 := createTestService("test-service", "127.0.0.1", 8081, 1.0)
		instances := []g.Service{service1, service2}
		ctx := context.Background()

		// 执行多次选择，统计选择结果
		service1Count := 0
		service2Count := 0
		totalRuns := 1000

		for i := 0; i < totalRuns; i++ {
			service, err := strategy.Select(ctx, instances)
			assert.NoError(t, err)

			if service == service1 {
				service1Count++
			} else if service == service2 {
				service2Count++
			}
		}

		// 验证权重效果，service1应该被选择的概率远高于service2
		assert.Greater(t, service1Count, service2Count*5, "权重为9的服务应该被选择的次数远多于权重为1的服务")
	})

	// 测试场景2: 服务实例列表为空
	t.Run("Empty instances list", func(t *testing.T) {
		// 创建空的实例列表
		instances := []g.Service{}
		ctx := context.Background()

		// 执行测试
		service, err := strategy.Select(ctx, instances)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "no available service instances")
	})

	// 测试场景3: 所有服务实例权重为0
	t.Run("All instances with zero weight", func(t *testing.T) {
		// 创建权重为0的服务实例
		service1 := createTestService("test-service", "127.0.0.1", 8080, 0.0)
		service2 := createTestService("test-service", "127.0.0.1", 8081, 0.0)
		instances := []g.Service{service1, service2}
		ctx := context.Background()

		// 执行测试
		service, err := strategy.Select(ctx, instances)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "total weight must be greater than zero")
	})
}
