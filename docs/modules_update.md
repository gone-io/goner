# Gone框架模块化改造之路

## 一、改造背景

### 1. Gone框架的定位

Gone框架是一个轻量级的依赖注入框架，它通过简洁的API和注解式依赖声明，帮助开发者轻松构建模块化、可测试的应用程序。作为一个轻量级框架，Gone坚持"小而美"的设计理念，只提供必要的核心功能，避免过度臃肿。

### 2. 可插拔组件的设计

Gone框架提供了一系列可插拔的组件（称为Goner），涵盖了Web开发、数据库访问、配置管理、日志记录等多个方面：

- gin：Web框架集成
- gorm：ORM框架集成
- viper：配置管理
- redis：缓存服务
- grpc：RPC服务
- apollo/nacos：配置中心
- zap：日志组件
- 等等

这种设计让开发者可以根据实际需求选择所需组件，实现灵活组合，避免了"全家桶"式的强制依赖。

### 3. 原有架构的问题

改造前，Gone核心代码和所有Goner组件都在同一个仓库中管理，共享同一个Go module。这种架构存在以下问题：

- **依赖膨胀**：用户即使只需要某一个特定功能（如gin集成），也必须导入整个Gone模块，包含大量不需要的依赖
- **版本管理困难**：所有组件共享同一个版本号，无法针对单个组件进行独立版本迭代
- **维护成本高**：随着组件数量增加，单一仓库的维护难度急剧上升
- **违背轻量级理念**：过多的依赖与Gone框架轻量级的设计理念相矛盾

## 二、改造过程

### 1. 仓库分离

首先，我们将Goner组件库从原始的`github.com/gone-io/gone`仓库中分离出来，独立为一个新仓库：`github.com/gone-io/goner`。

这一步骤包括：

- 创建新的Git仓库
- 迁移相关代码，保留提交历史
- 调整导入路径

### 2. 单仓库多模块（Mono Repo）改造

对`github.com/gone-io/goner`进行Mono Repo改造，将每个Goner组件作为独立的子模块管理：

```
github.com/gone-io/goner/
├── g/                  # 基础接口定义
├── gin/                # Gin框架集成
├── gorm/               # GORM框架集成
├── redis/              # Redis客户端集成
├── viper/              # Viper配置管理
├── apollo/             # Apollo配置中心
├── nacos/              # Nacos配置中心
├── zap/                # Zap日志组件
├── grpc/               # gRPC服务集成
├── cmux/               # 连接复用
├── schedule/           # 定时任务
├── tracer/             # 链路追踪
├── urllib/             # HTTP客户端
├── xorm/               # XORM框架集成
└── ...
```

每个子目录都是一个独立的Go模块，拥有自己的`go.mod`文件，可以独立版本化和发布。

### 3. 接口抽象与依赖解耦

为了解决组件间的相互依赖问题，我们采用了"依赖接口，不依赖实现"的设计原则：

- 创建`github.com/gone-io/goner/g`子模块，用于定义所有公共接口
- 所有组件都依赖于这个基础接口模块，而不是直接依赖其他组件
- 通过接口抽象，实现了组件间的松耦合

例如，在`g/tracer.go`中定义了链路追踪相关接口，而不是直接依赖于`tracer`的具体实现：

```go
// /g/tracer.go
package g

// Tracer 用于日志追踪，为同一调用链分配统一的traceId，方便日志追踪
type Tracer interface {

    //SetTraceId 为调用函数设置traceId。如果traceId为空字符串，将自动生成一个。
    //调用函数可以通过GetTraceId()方法获取traceId。
    SetTraceId(traceId string, fn func())

    //GetTraceId 获取当前goroutine的traceId
    GetTraceId() string
 
    //Go 启动一个新的goroutine，代替原生的`go func`，可以将traceId传递给新的goroutine
    Go(fn func())
}
```

## 三、存在问题及解决方案

### 1. Go项目单仓库多模块的版本管理

在Go模块系统中，单仓库多模块的版本管理是一个常见挑战。我们采用了以下策略：

#### 子模块版本标签

我们为每个子模块设置独立的Git标签进行版本管理，确保依赖引用时能准确定位对应版本的代码：

```
<module>/<version>
```

例如：
- `gin/v1.0.0` - 表示gin模块的1.0.0版本
- `gorm/v0.2.1` - 表示gorm模块的0.2.1版本

这种方式允许每个子模块独立版本化，同时保持在同一个仓库中管理的便利性。

#### 版本依赖管理

在子模块的`go.mod`文件中，通过`replace`指令处理本地模块间的依赖关系：

```go
// 开发环境中使用本地路径
replace github.com/gone-io/goner/g => ../g

// 发布后，使用具体版本
require github.com/gone-io/goner/g v0.1.0
```

这确保了在开发过程中可以使用本地模块，而在发布后则使用特定版本的依赖。

### 2. Go项目单仓库多模块的单元测试

对于单仓库多模块中的单元测试，我们采用了分散测试、集中报告的策略：

#### 分模块测试

每个子模块都有自己的测试套件，可以独立运行和验证：

```bash
cd gin && go test -v ./...
cd gorm && go test -v ./...
```

#### 集中测试报告

通过GitHub Actions工作流（`.github/workflows/go.yml`）实现测试的自动化运行和报告合并：

```yaml
- name: Run coverage
  run: find . -name go.mod -not -path "*/example/*" -not -path "*/examples/*" | xargs -n1 dirname | xargs -L1 bash -c 'cd "$0" && pwd && go test -race -coverprofile=coverage.txt -covermode=atomic ./... || exit 255' || echo "Tests failed"
- name: Merge coverage
  run: find . -name "coverage.txt" -exec cat {} \; > total_coverage.txt
- name: Upload coverage reports to Codecov
  uses: codecov/codecov-action@v5
  with:
      token: ${{ secrets.CODECOV_TOKEN }}
      files: ./total_coverage.txt
```

这个工作流会：
1. 查找所有的`go.mod`文件（排除示例目录）
2. 在每个模块目录中运行测试并生成覆盖率报告
3. 合并所有覆盖率报告
4. 将合并后的报告上传到Codecov进行分析

### 3. 组件间依赖问题

Goner组件之间存在相互依赖关系，这可能导致循环依赖或过度耦合。我们通过以下方式解决：

#### 接口抽象

将所有共享接口抽象到`github.com/gone-io/goner/g`子模块中：

```go
// g/cmux.go
package g

// Listener 定义了网络监听器接口
type Listener interface {
    // ...
}

// g/tracer.go
package g

// Tracer 定义了链路追踪接口
type Tracer interface {
    // ...
}
```

#### 依赖注入

利用Gone框架的依赖注入机制，在运行时动态解析依赖关系：

```go
type MyService struct {
    gone.Flag
    Logger g.Logger `gone:"*"`  // 依赖接口，而非具体实现
}
```

#### 插件化设计

采用插件化设计，允许组件在运行时动态加载：

```go
func main() {
    gone.
        Loads(
            zap.Load,     // 加载zap日志组件
            gin.Load,     // 加载gin组件
            gorm.Load,    // 加载gorm组件
        ).
        Run(func() {
            // 应用启动
        })
}
```

## 四、改造成果

通过这次模块化改造，Gone框架及其生态系统获得了显著改进：

1. **更轻量的依赖**：用户可以只引入所需的组件，大大减少不必要的依赖
2. **灵活的版本管理**：每个组件可以独立版本化，便于针对性升级和维护
3. **更好的可测试性**：组件间的松耦合使单元测试更易编写和维护
4. **简化的开发流程**：明确的接口定义和依赖关系，使开发新组件变得更加简单
5. **更符合Go设计理念**：遵循"小而美"的设计原则，每个模块专注于解决特定问题

## 五、未来展望

Gone框架的模块化改造是一个持续的过程，未来我们计划：

1. **完善文档体系**：为每个组件提供详细的使用文档和示例
2. **扩展组件生态**：开发更多实用组件，如消息队列、分布式锁等
3. **性能优化**：针对核心组件进行性能优化，提高整体效率
4. **社区建设**：鼓励社区贡献，建立更加活跃的开源生态

通过这些努力，我们希望Gone框架能够成为Go语言生态中一个轻量、高效、易用的依赖注入框架，为开发者提供更好的开发体验。

## 参考资源

- [Gone框架官方仓库](https://github.com/gone-io/gone)
- [Goner组件库](https://github.com/gone-io/goner)
- [Go Modules文档](https://go.dev/ref/mod)