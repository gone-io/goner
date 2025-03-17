package tracer

import (
	"testing"

	"github.com/gone-io/gone/v2"
)

// 模拟gone.Loader接口
type mockLoader struct {
	loadedComponents []interface{}
}

func (m *mockLoader) Load(component gone.Goner, opts ...gone.Option) error {
	m.loadedComponents = append(m.loadedComponents, component)
	return nil
}

func (m *mockLoader) Loaded(gone.LoaderKey) bool {
	return false
}

func (m *mockLoader) GetComponents() []interface{} {
	return m.loadedComponents
}

func TestLoad(t *testing.T) {
	// 创建模拟加载器
	mockLoader := &mockLoader{}

	// 调用Load函数
	err := Load(mockLoader)

	// 验证没有错误
	if err != nil {
		t.Errorf("Load() returned error: %v", err)
	}

	// 验证加载了正确的组件
	if len(mockLoader.loadedComponents) != 1 {
		t.Errorf("Expected 1 component to be loaded, got %d", len(mockLoader.loadedComponents))
	}

	// 验证加载的组件类型是否正确
	if _, ok := mockLoader.loadedComponents[0].(*tracer); !ok {
		t.Errorf("Expected component of type *tracer, got %T", mockLoader.loadedComponents[0])
	}
}

func TestLoadGidTracer(t *testing.T) {
	// 创建模拟加载器
	mockLoader := &mockLoader{}

	// 调用LoadGidTracer函数
	err := LoadGidTracer(mockLoader)

	// 验证没有错误
	if err != nil {
		t.Errorf("LoadGidTracer() returned error: %v", err)
	}

	// 验证加载了正确的组件
	if len(mockLoader.loadedComponents) != 1 {
		t.Errorf("Expected 1 component to be loaded, got %d", len(mockLoader.loadedComponents))
	}

	// 验证加载的组件类型是否正确
	if _, ok := mockLoader.loadedComponents[0].(*tracerOverGid); !ok {
		t.Errorf("Expected component of type *tracerOverGid, got %T", mockLoader.loadedComponents[0])
	}
}
