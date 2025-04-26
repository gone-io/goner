# Gone Simple 示例

## 项目概述

这是一个 Gone 框架的最简单示例项目，展示了 Gone 框架的基本使用方法和依赖注入特性。通过这个示例，你可以快速了解：

- Gone 框架的基本项目结构
- 依赖注入的使用方式
- 配置文件的读取方式
- 接口与实现的分离原则

## 项目结构

```
.
├── cmd/                # 程序入口目录
│   ├── import.gone.go  # 依赖导入文件
│   └── main.go        # 主程序入口
├── config/            # 配置文件目录
│   └── default.yaml   # 默认配置文件
├── implement/         # 接口实现目录
│   ├── hello.go      # Service 接口实现
│   └── init.gone.go  # 实现初始化文件
├── service/          # 接口定义目录
│   └── interface.go  # Service 接口定义
└── module.load.go    # 模块加载文件
```

## 功能特点

1. **简洁的依赖注入**：通过 `gone:"*"` 标签实现接口注入
2. **配置自动绑定**：使用 `gone:"config,app.name"` 自动绑定配置项
3. **接口与实现分离**：遵循 Go 语言的最佳实践，将接口定义和实现分离
4. **模块化结构**：清晰的目录结构，便于项目扩展

## 使用说明

### 1. 安装依赖

首先安装 Gone 框架的配置组件：

```bash
gonectr install github.com/gone-io/goner/viper
```

### 2. 运行项目

在项目根目录下执行：

```bash
gonectr run ./cmd
```

### 3. 预期输出

程序将输出：
```
hello simple-app
```

## 核心代码说明

### 接口定义

```go
// service/interface.go
type Service interface {
    SayHello(name string) string
}
```

### 接口实现

```go
// implement/hello.go
type serviceImpl struct {
    gone.Flag
}

func (s serviceImpl) SayHello(name string) string {
    return fmt.Sprintf("hello %s", name)
}
```

### 主程序

```go
// cmd/main.go
func main() {
    gone.Run(func(in struct {
        service service.Service `gone:"*"`
        appName string          `gone:"config,app.name"`
    }) {
        println(in.service.SayHello(in.appName))
    })
}
```

## 配置说明

配置文件位于 `config/default.yaml`：

```yaml
app:
  name: simple-app
```

## 扩展建议

1. 可以添加更多的配置项来展示配置绑定能力
2. 可以增加多个接口实现来展示依赖注入的灵活性
3. 可以添加数据库访问等实际业务场景的示例

## 相关文档

- [Gone 框架文档](https://github.com/gone-io/gone)
- [配置中心示例](https://github.com/gone-io/goner/tree/main/examples/config_center)
- [更多示例项目](https://github.com/gone-io/goner/tree/main/examples)