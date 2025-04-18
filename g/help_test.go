package g

import (
	"errors"
	"github.com/gone-io/gone/v2"
	"reflect"
	"testing"

	mock "github.com/gone-io/gone/mock/v2"
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

func TestBuildOnceLoadFunc(t *testing.T) {
	type gTest struct {
		gone.Flag
	}

	loadFunc := BuildOnceLoadFunc(L(&gTest{}), F(func(loader gone.Loader) error {
		return nil
	}))
	gone.NewApp(loadFunc, loadFunc).Run(func(
		gList []*gTest,
	) {
		assert.Len(t, gList, 1)
	})
}

func TestBuildOnceLoadFuncError(t *testing.T) {
	type gTest struct {
		gone.Flag
	}

	loadFunc := BuildOnceLoadFunc(
		L(&gTest{}, gone.Name("test")),
		L(&gTest{}, gone.Name("test")),
	)

	assert.Panics(t, func() {
		gone.
			NewApp(loadFunc, loadFunc).
			Run(func(
				gList []*gTest,
			) {
			})
	})
}

func TestBuildOnceLoadFuncError2(t *testing.T) {
	type gTest struct {
		gone.Flag
	}

	loadFunc := BuildOnceLoadFunc(
		F(func(loader gone.Loader) error {
			return errors.New("test")
		}),
	)

	assert.Panics(t, func() {
		gone.
			NewApp(loadFunc, loadFunc).
			Run(func(
				gList []*gTest,
			) {
			})
	})
}

func TestSingLoadProviderFunc(t *testing.T) {
	t.Run("once load", func(t *testing.T) {
		type TestStruct struct {
		}

		var loadTimes int

		providerFunc := SingLoadProviderFunc(func(tagConf string, param struct{}) (*TestStruct, error) {
			loadTimes++
			return &TestStruct{}, nil
		})
		gone.
			NewApp(providerFunc, providerFunc).
			Run(func(s *TestStruct, in struct {
				s1 *TestStruct
				s2 *TestStruct
			}) {
				assert.Equal(t, 1, loadTimes)
			})
	})
}

func TestNamedThirdComponentLoadFunc(t *testing.T) {
	t.Run("once load", func(t *testing.T) {
		type TestStruct struct {
			Name string
		}
		var s = TestStruct{
			Name: "X",
		}

		loadFunc := NamedThirdComponentLoadFunc("test", &s)
		gone.
			NewApp(loadFunc, loadFunc).
			Run(func(s0 *TestStruct, in struct {
				s1 *TestStruct `gone:"test"`
				s2 *TestStruct `gone:"test"`
			}) {
				assert.Equal(t, s0, &s)
				assert.Equal(t, s0, in.s1)
				assert.Equal(t, s0, in.s2)
			})
	})
}

type provided struct {
	gone.Flag
}

type provider struct {
	gone.Flag
}

func (s *provider) Provide() (*provided, error) {
	return &provided{}, nil
}

type provider2 struct {
	gone.Flag
}

func (s *provider2) Provide(tagConf string) (*provided, error) {
	return &provided{}, nil
}

type provider3 struct {
	gone.Flag
	err error
	c   any
}

func (s *provider3) GonerName() string {
	return "test"
}
func (s *provider3) Provide(tagConf string, t reflect.Type) (any, error) {
	return s.c, s.err
}

func TestGetComponentByName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	keeper := mock.NewMockGonerKeeper(controller)

	t.Run("not found", func(t *testing.T) {
		keeper.EXPECT().GetGonerByName("test").Return(nil)
		component, err := GetComponentByName[*provided](keeper, "test")
		assert.Error(t, err)
		assert.Nil(t, component)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("found component by name", func(t *testing.T) {
		keeper.EXPECT().GetGonerByName("test").Return(&provided{})
		component, err := GetComponentByName[*provided](keeper, "test")
		assert.Nil(t, err)
		assert.NotNil(t, component)
	})

	t.Run("found NoneParamProvider", func(t *testing.T) {
		keeper.EXPECT().GetGonerByName("test").Return(&provider{})
		component, err := GetComponentByName[*provided](keeper, "test")
		assert.Nil(t, err)
		assert.NotNil(t, component)
	})

	t.Run("found Provider", func(t *testing.T) {
		keeper.EXPECT().GetGonerByName("test").Return(&provider2{})

		component, err := GetComponentByName[*provided](keeper, "test")
		assert.Nil(t, err)
		assert.NotNil(t, component)
	})

	t.Run("found NamedProvider", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			keeper.EXPECT().GetGonerByName("test").Return(&provider3{c: &provided{}})
			component, err := GetComponentByName[*provided](keeper, "test")
			assert.Nil(t, err)
			assert.NotNil(t, component)
		})
		t.Run("error", func(t *testing.T) {
			keeper.EXPECT().GetGonerByName("test").Return(&provider3{err: errors.New("test"), c: &provided{}})
			component, err := GetComponentByName[*provided](keeper, "test")
			assert.Error(t, err)
			assert.Nil(t, component)
		})
	})
	t.Run("found other type", func(t *testing.T) {
		keeper.EXPECT().GetGonerByName("test").Return("test")
		component, err := GetComponentByName[*provided](keeper, "test")
		assert.Nil(t, component)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not found compatible component")
	})
}
