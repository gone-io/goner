[//]: # (desc: MCP 服务 快速入门示例)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone MCP 快速入门示例

这是一个基于 Gone MCP 组件的快速入门示例，展示了如何创建一个简单的 MCP 工具服务器和客户端。通过这个示例，你可以快速了解和掌握 Gone MCP 组件的基本使用方法。

## 项目结构

```
.
├── client/         # 客户端示例代码
│   └── main.go     # 客户端主程序
├── go.mod          # Go 模块定义
└── server/         # 服务端示例代码
    └── main.go     # 服务端主程序
```

## 功能说明

这个示例实现了一个简单的问候服务：

1. **服务端**：
   - 实现了一个名为 `hello_world` 的 MCP 工具
   - 接受一个必填的字符串参数 `name`
   - 返回格式化的问候语

2. **客户端**：
   - 通过 stdio 方式与服务端通信
   - 展示了完整的 MCP 客户端使用流程：
     - 初始化连接
     - 获取可用工具列表
     - 调用工具并处理结果

## 使用方法

### 前置条件

确保已安装 Go 1.16 或更高版本。

### 运行示例

1. 克隆项目并进入示例目录：
   ```bash
   cd examples/mcp/quick_start
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```

3. 运行客户端程序：
   ```bash
   go run ./client
   ```

   客户端程序会自动启动服务端，并执行以下操作：
   - 初始化与服务端的连接
   - 获取并显示可用工具列表
   - 调用 `hello_world` 工具并显示结果

## 示例输出

运行客户端后，你将看到类似以下的输出：

```
Initialized with server: quick-start 0.0.1

Listing available tools...
- hello_world: Say hello to someone

Calling `hello_world`
Hello, John!
```

## 代码说明

### 服务端 (server/main.go)

- 使用 `gone.Flag` 实现 Goner 定义
- 通过 `Define()` 方法定义工具名称、描述和参数
- 在 `Handler()` 方法中实现工具的具体逻辑

### 客户端 (client/main.go)

- 使用依赖注入方式获取 MCP 客户端
- 演示了完整的客户端初始化和工具调用流程
- 包含了错误处理和结果展示的最佳实践

## 扩展建议

1. 尝试添加更多参数到 `hello_world` 工具
2. 实现新的工具，如计算器或文本处理工具
3. 尝试使用其他通信方式，如 HTTP 或 WebSocket

## 相关文档

- [Gone MCP 组件文档](../../../mcp)
- [Gone 框架文档](https://github.com/gone-io/gone)