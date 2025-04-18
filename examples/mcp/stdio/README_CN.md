# Gone MCP Stdio 示例

这是一个基于 Gone MCP 组件的 stdio 通信示例，展示了如何使用标准输入输出（stdio）方式创建 MCP 工具服务器和客户端。通过这个示例，你可以了解如何在需要进程间通信的场景中使用 Gone MCP 组件。

## 项目结构

```
.
├── client/         # 客户端示例代码
│   └── main.go     # 客户端主程序
├── config/         # 配置文件目录
│   └── default.yaml # 默认配置文件
├── go.mod          # Go 模块定义
└── server/         # 服务端示例代码
    ├── functional_add/  # 功能模块目录
    ├── goner_define/    # Gone 定义文件
    ├── import.gone.go   # Gone 导入文件
    └── main.go         # 服务端主程序
```

## 功能说明

这个示例展示了如何使用 stdio 通信方式实现 MCP 服务：

1. **服务端**：
   - 实现了基于 stdio 的 MCP 工具服务器
   - 支持标准输入输出流通信
   - 包含完整的 Gone 项目结构

2. **客户端**：
   - 通过 stdio 方式与服务端通信
   - 展示了完整的 MCP 客户端使用流程：
     - 初始化连接
     - 获取可用工具列表
     - 调用工具并处理结果

## 使用场景

stdio 通信方式特别适合以下场景：

1. **进程间通信**：当需要在同一台机器上的不同进程之间进行通信时
2. **命令行工具**：开发需要与其他程序交互的命令行工具
3. **插件系统**：实现主程序和插件之间的通信
4. **调试环境**：在开发和测试阶段，方便查看通信数据

## 使用方法

### 前置条件

确保已安装 Go 1.16 或更高版本。

### 运行示例

1. 进入示例目录：
   ```bash
   cd examples/mcp/stdio
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```
3. 生成辅助代码
   ```bash
   go generator./...
   ```

4. 运行客户端程序：
   ```bash
   go run ./client
   ```

   客户端程序会自动启动服务端，并通过 stdio 方式进行通信。

## 配置说明

### 服务端配置

服务端配置位于 `config/default.yaml`，主要包含：

- 服务器基本信息
- 工具配置
- 日志配置

### 客户端配置

客户端通过代码配置，主要设置：

```go
client *client.Client `gone:"*,type=stdio,param=go run ./server"`
```

- `type=stdio`：指定使用 stdio 通信方式
- `param=go run ./server`：指定服务端启动命令

## 扩展建议

1. 添加更多的工具实现
2. 实现双向通信功能
3. 添加数据传输的压缩和加密
4. 实现更复杂的进程间通信场景

## 注意事项

1. stdio 通信方式仅适用于本地进程间通信
2. 确保服务端和客户端使用相同的协议版本
3. 注意处理 stdio 流的缓冲和关闭
4. 建议在开发环境中使用，生产环境可能需要考虑其他通信方式

## 相关文档

- [Gone MCP 组件文档](../../../mcp)
- [Gone 框架文档](https://github.com/gone-io/gone)
- [Gone Viper 读取配置文档](../../../viper)
- [mcp-go](github.com/mark3labs/mcp-go)