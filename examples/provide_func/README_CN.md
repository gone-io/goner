[//]: # (desc: Provide函数示例对接第三方组件)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone框架Provide函数示例

本示例展示了如何使用 Gone 框架的 **Provide函数** 来优雅地集成第三方组件。通过简单的标签配置和依赖注入机制，您可以轻松地将第三方组件引入到 Gone 框架中，实现组件的灵活配置和高效管理。

## 核心特性

- 标准化的 **Provide函数** 实现模式
- 强大的标签配置解析机制
- 自动依赖注入和日志管理
- 专为第三方组件集成优化

## 技术实现

### 组件定义

```go
// ThirdComponent 模拟第三方组件
type ThirdComponent struct {
}
```

### Provide函数实现

```go
func provide(tagConf string, in struct {
	logger gone.Logger `gone:"*"`
}) (*ThirdComponent, error) {
	confMap, confKeys := gone.TagStringParse(tagConf)
	in.logger.Infof("confMap => %#v\nconfKeys=>%#v", confMap, confKeys)

	// 根据不同配置创建第三方组件
	return &ThirdComponent{}, nil
}
```

**Provide函数参数说明：**
- `tagConf string`：接收组件标签配置，支持 key=value 格式和位置参数
- `in struct`：依赖注入结构体，用于注入所需的系统组件
  - `logger gone.Logger`：系统日志组件，通过 `gone:"*"` 标签注入

**返回值说明：**
- `*ThirdComponent`：返回创建的组件实例
- `error`：返回可能的错误信息

### 组件注册

```go
func Load(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(provide))
}
```

通过 `gone.WrapFunctionProvider` 将 Provide 函数包装为标准的 Gone 框架组件。

## 使用指南

### 基础用法

```go
type YourComponent struct {
    third *ThirdComponent `gone:"*"`
}
```

### 配置驱动

```go
type YourComponent struct {
    third *ThirdComponent `gone:"*,key=value,another=123"`
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
- [依赖注入最佳实践指南](../../docs/wrap-function-provider.md)