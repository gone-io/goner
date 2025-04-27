[//]: # (desc: 定义Provider组件对接第三方)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone框架标准Provider组件示例

本示例展示了如何使用 Gone 框架的 **Provider接口** 来优雅地集成第三方组件。通过实现标准的Provider接口，您可以轻松地将第三方组件引入到 Gone 框架中，实现组件的灵活配置和高效管理。

## 核心特性

- 标准化的 **Provider接口** 实现模式
- 支持带参数和无参数的Provider实现
- 自动依赖注入和组件管理
- 专为第三方组件集成优化

## 技术实现

### 组件定义

```go
// ThirdComponent1 模拟第三方组件
type ThirdComponent1 struct {
}

// ThirdComponent2 模拟第三方组件
type ThirdComponent2 struct {
}
```

### Provider接口实现

#### 带参数的Provider实现

```go
type provider struct {
    gone.Flag
}

func (p *provider) Provide(tagConf string) (*ThirdComponent1, error) {
    // 根据配置创建第三方组件
    return &ThirdComponent1{}, nil
}
```

**Provider接口说明：**
- `Provide(tagConf string)`：接收组件标签配置，支持 key=value 格式和位置参数
- 返回值：返回创建的组件实例和可能的错误信息

#### 无参数的Provider实现

```go
type noneParamProvider struct {
    gone.Flag
}

func (p noneParamProvider) Provide() (*ThirdComponent2, error) {
    // 创建第三方组件
    return &ThirdComponent2{}, nil
}
```

### 组件注册

```go
func Load(loader gone.Loader) error {
    loader.
        MustLoad(&provider{}).
        MustLoad(&noneParamProvider{})
    return nil
}
```

## 使用指南

### 基础用法

```go
type YourComponent struct {
    comp1 *ThirdComponent1 `gone:"*"`
    comp2 *ThirdComponent2 `gone:"*"`
}
```

### 配置驱动

```go
type YourComponent struct {
    comp1 *ThirdComponent1 `gone:"*,key=value,another=123"`
    comp2 *ThirdComponent2 `gone:"*"`
}
```

**标签配置示例：**
- 基础注入：`gone:"*"`
- 带配置注入：`gone:"*,name=myComponent"`
- 多配置项：`gone:"*,host=localhost,port=8080,timeout=5s"`

## 最佳实践

### 适用场景

- 需要动态配置的组件初始化
- 第三方组件的标准化集成
- 复杂组件初始化逻辑的封装
- 组件生命周期的精确管理

### 配置管理

- 使用 `gone.TagStringParse` 解析标签配置
  - 返回 `confMap`：包含所有 key=value 形式的配置
  - 返回 `confKeys`：保持配置项的原始顺序
- 支持默认值和必填项校验
- 配置项支持多种格式：
  - 位置参数：`gone:"*,redis"`
  - 键值对：`gone:"*,driver=redis"`
  - 混合使用：`gone:"*,redis,port=6379"`

### 错误处理

- 配置校验
  - 检查必填配置项
  - 验证配置值的格式和范围
  - 处理配置冲突

- 组件初始化
  - 捕获并包装第三方组件的错误
  - 提供清晰的错误信息
  - 实现错误恢复机制

## 扩展资源

- [Gone框架官方文档](https://github.com/gone-io/gone)
- [Provider 机制介绍](https://github.com/gone-io/gone/blob/main/docs/provider.md)
