package g

import (
	"testing"

	mock "github.com/gone-io/gone/mock/v2"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestGetLocalIps 测试成功获取本地IP地址
func TestGetLocalIps(t *testing.T) {
	// 执行测试
	ips := GetLocalIps()

	// 验证结果
	assert.NotNil(t, ips)
	for _, ip := range ips {
		// 验证IP地址格式是否正确
		assert.NotNil(t, ip.To4())
		// 验证不是环回地址
		assert.False(t, ip.IsLoopback())
	}
}

// TestRecover 测试Recover函数的panic恢复和日志记录
func TestRecover(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLogger
	mockLogger := mock.NewMockLogger(ctrl)

	// 设置期望：记录错误日志
	mockLogger.EXPECT().
		Errorf(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	// 执行测试：触发panic并恢复
	func() {
		defer Recover(mockLogger)
		panic("test panic")
	}()
}

// TestBuildLoadFunc_Success 测试成功加载场景
func TestBuildLoadFunc_Success(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader和MockGoner
	mockLoader := mock.NewMockLoader(ctrl)
	mockGoner := &struct {
		gone.Flag
	}{}

	// 设置期望：检查是否已加载
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：成功加载Goner
	mockLoader.EXPECT().
		Load(mockGoner, gomock.Any()).
		Return(nil)

	// 执行测试
	err := BuildLoadFunc(mockLoader, L(mockGoner))

	// 验证结果
	assert.NoError(t, err)
}

// TestBuildLoadFunc_LoadError 测试加载失败场景
func TestBuildLoadFunc_LoadError(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader和MockGoner
	mockLoader := mock.NewMockLoader(ctrl)
	mockGoner := &struct {
		gone.Flag
	}{}

	// 设置期望：检查是否已加载
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

	// 设置期望：加载失败
	expectedErr := gone.NewError(400, "mock load error", 500)
	mockLoader.EXPECT().
		Load(mockGoner, gomock.Any()).
		Return(expectedErr)

	// 执行测试
	err := BuildLoadFunc(mockLoader, L(mockGoner))

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestBuildLoadFunc_AlreadyLoaded 测试已经加载过的场景
func TestBuildLoadFunc_AlreadyLoaded(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockLoader
	mockLoader := mock.NewMockLoader(ctrl)

	// 设置期望：检查是否已加载，返回true表示已加载
	mockLoader.EXPECT().Loaded(gomock.Any()).Return(true)

	// 执行测试
	err := BuildLoadFunc(mockLoader)

	// 验证结果
	assert.NoError(t, err)
}
