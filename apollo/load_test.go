package apollo

import (
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
)

// 模拟gone.Loader接口
type mockLoader struct {
	loadedComponents []interface{}
	loadOptions      [][]gone.Option
	loadError        error
}

func (m *mockLoader) Load(component gone.Goner, opts ...gone.Option) error {
	m.loadedComponents = append(m.loadedComponents, component)
	m.loadOptions = append(m.loadOptions, opts)
	return m.loadError
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
	assert.NoError(t, err)

	// 验证加载了正确的组件
	assert.Equal(t, 2, len(mockLoader.loadedComponents), "应该加载了2个组件")

	// 验证第一个加载的组件是apolloClient
	assert.IsType(t, &apolloConfigure{}, mockLoader.loadedComponents[0], "第一个组件应该是apolloClient")

	// 验证第二个加载的组件是changeListener
	assert.IsType(t, &changeListener{}, mockLoader.loadedComponents[1], "第二个组件应该是changeListener")

	// 验证apolloClient加载时使用了正确的选项
	assert.Equal(t, 3, len(mockLoader.loadOptions[0]), "apolloClient应该有4个加载选项")
}

func TestLoadError(t *testing.T) {
	// 创建模拟加载器，设置加载错误
	mockLoader := &mockLoader{loadError: gone.NewError(1100, "模拟加载错误", 400)}

	// 调用Load函数
	err := Load(mockLoader)

	// 验证返回了错误
	assert.Error(t, err)
	assert.Equal(t, 1100, err.(gone.Error).Code())

	// 验证只尝试加载了第一个组件
	assert.Equal(t, 1, len(mockLoader.loadedComponents), "应该只尝试加载了1个组件")
}
