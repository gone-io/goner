package gorm

import (
	mock "github.com/gone-io/gone/mock/v2"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestLoad_Success 测试正常加载场景
func TestLoad_Success(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)

	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：成功加载iLogger
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&iLogger{}), gomock.Any()).
		Return(nil)

	// 设置期望：成功加载dbProvider
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&dbProvider{}), gomock.Any()).
		Return(nil)

	// 执行测试
	err := Load(mockLoader)

	// 验证结果
	assert.NoError(t, err)
}

// TestLoad_LoggerError 测试加载iLogger失败的场景
func TestLoad_LoggerError(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：加载iLogger失败
	expectedErr := gone.NewError(400, "mock logger error", 500)
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&iLogger{}), gomock.Any()).
		Return(expectedErr)

	// 执行测试
	err := Load(mockLoader)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestLoad_ProviderError 测试加载dbProvider失败的场景
func TestLoad_ProviderError(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：成功加载iLogger
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&iLogger{}), gomock.Any()).
		Return(nil)

	// 设置期望：加载dbProvider失败
	expectedErr := gone.NewError(400, "mock provider error", 500)
	mockLoader.EXPECT().
		Load(gomock.AssignableToTypeOf(&dbProvider{}), gomock.Any()).
		Return(expectedErr)

	// 执行测试
	err := Load(mockLoader)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
