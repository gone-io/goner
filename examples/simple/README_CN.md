[//]: # (desc: 简单示例，使用viper提供配置读取)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone Simple 示例

## 项目概述

这是一个 Gone 框架的简单示例项目，展示了 Gone 框架的基本使用方法和依赖注入特性。通过这个示例，你可以快速了解：

- Gone 框架的标准项目结构
- 依赖注入的使用方式
- 配置文件的读取方式（使用 Viper）
- 接口与实现的分离原则
- CLI 和 Server 两种应用模式

## 项目结构

```
.
├── cmd/                    # 程序入口目录
│   ├── cli/               # 命令行应用入口
│   │   ├── import.gone.go # 依赖导入文件
│   │   └── main.go        # CLI 主程序
│   └── server/            # 服务器应用入口
│       ├── import.gone.go # 依赖导入文件
│       ├── init.gone.go   # 初始化文件
│       └── main.go        # Server 主程序
├── config/                # 配置文件目录
│   ├── default.yaml       # 默认配置文件
│   ├── dev.yaml          # 开发环境配置
│   ├── local.yaml        # 本地环境配置
│   ├── prod.yaml         # 生产环境配置
│   └── test.yaml         # 测试环境配置
├── internal/              # 内部代码目录
│   ├── controller/        # 控制器目录
│   │   └── hello.go      # Hello 控制器
│   ├── interface/         # 接口定义目录
│   │   ├── entity/       # 实体定义
│   │   ├── mock/         # Mock 实现
│   │   └── service/      # 服务接口
│   │       └── i_server.go # IService 接口定义
│   ├── module/           # 模块实现目录
│   │   ├── hello/        # Hello 模块
│   │   │   └── hello.go  # Service 接口实现
│   │   └── user/         # User 模块
│   ├── pkg/              # 内部工具包
│   │   ├── e/           # 错误定义
│   │   └── utils/       # 工具函数
│   └── router/           # 路由定义
│       ├── auth_router.go # 认证路由
│       └── pub_router.go  # 公共路由
├── asserts/              # 静态资源目录
├── docker-compose.yaml   # Docker Compose 配置
├── Dockerfile           # Docker 镜像构建文件
├── Makefile            # Make 构建脚本
├── module.load.go      # 模块加载文件
└── pacakge.go          # 包定义文件
```

## 功能特点

1. **简洁的依赖注入**：通过 `gone:"*"` 标签实现接口注入
2. **配置自动绑定**：使用 `gone:"config,app.name"` 自动绑定配置项
3. **接口与实现分离**：遵循 Go 语言的最佳实践，将接口定义和实现分离
4. **模块化结构**：采用 internal 目录结构，清晰的模块划分
5. **多环境配置**：支持不同环境的配置文件
6. **双模式支持**：同时支持 CLI 和 Server 两种运行模式

## 使用说明

### 1. 安装依赖

项目已经配置了必要的依赖，包括：
- `github.com/gone-io/gone/v2` - Gone 核心框架
- `github.com/gone-io/goner/viper` - Viper 配置组件
- `github.com/gone-io/goner/gin` - Gin Web 框架组件

### 2. 运行 CLI 应用

在项目根目录下执行：

```bash
go run ./cmd/cli
```

### 3. 运行 Server 应用
如果需要使用gin web框架，在项目根目录下执行：
```bash
gonectl install goner/gin Loader
```

在项目根目录下执行：

```bash
go run ./cmd/server
```

### 4. 预期输出

**CLI 应用输出：**
```
hello root-app
```

**Server 应用输出：**
```
after server start
press `ctr + c` to stop!
```

## 核心代码说明

### 接口定义

```go
// internal/interface/service/i_server.go
type IService interface {
    SayHello(name string) string
}
```

### 接口实现

```go
// internal/module/hello/hello.go
type serviceImpl struct {
    gone.Flag
}

func (s serviceImpl) SayHello(name string) string {
    return fmt.Sprintf("hello %s", name)
}
```

### CLI 主程序

```go
// cmd/cli/main.go
func main() {
    gone.Run(func(in struct {
        service service.IService `gone:"*"`
        appName string           `gone:"config,app.name"`
    }) {
        println(in.service.SayHello(in.appName))
    })
}
```

### Server 主程序

```go
// cmd/server/main.go
func main() {
    gone.Serve()
}
```

### 模块加载配置

```go
// module.load.go
var loaders = []gone.LoadFunc{
    viper.Load,
}

func GoneModuleLoad(loader gone.Loader) error {
    var ops []*g.LoadOp
    for _, f := range loaders {
        ops = append(ops, g.F(f))
    }
    return g.BuildOnceLoadFunc(ops...)(loader)
}
```

## 配置说明

配置文件位于 `config/default.yaml`：

```yaml
app:
  name: root-app
  env: dev
```

项目支持多环境配置，可以通过环境变量或启动参数指定不同的配置文件。

## Docker 支持

项目包含了 Docker 相关配置：

- `Dockerfile` - 用于构建应用镜像
- `docker-compose.yaml` - 用于本地开发和测试
- `.dockerignore` - Docker 构建时忽略的文件

## 扩展建议

1. **添加 Web 接口**：可以在 controller 目录下添加 HTTP 接口
2. **数据库集成**：可以集成 GORM 或 XORM 进行数据库操作
3. **中间件支持**：可以添加认证、日志、限流等中间件
4. **配置扩展**：可以添加更多的配置项来展示配置绑定能力
5. **测试覆盖**：可以添加单元测试和集成测试

## 相关文档

- [Gone 框架文档](https://github.com/gone-io/gone)
- [Goner 组件库](https://github.com/gone-io/goner)
- [更多示例项目](https://github.com/gone-io/goner/tree/main/examples)