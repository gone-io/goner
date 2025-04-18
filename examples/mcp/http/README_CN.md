# Gone MCP HTTP Demo

这个示例展示了如何使用 Gone MCP 组件构建基于 HTTP 的客户端和服务端应用。

## 项目结构

```
.
├── client/         # 客户端示例代码
│   └── main.go     # 客户端主程序
├── config/         # 配置文件目录
│   └── default.yaml # 默认配置文件
├── go.mod          # Go 模块定义
└── server/         # 服务端示例代码
    ├── functional_add/    # 功能定义
    ├── goner_define/      # Gone 组件定义
    ├── import.gone.go     # Gone 导入文件
    └── main.go            # 服务端主程序
```

## 功能说明

这个示例实现了以下功能：

1. **资源服务**：提供用户资源访问
   - 支持通过 URI 模板 `users://{id}/profile` 访问用户信息
   - 返回 JSON 格式的用户数据

2. **提示服务**：代码审查功能
   - 提供代码审查辅助功能
   - 接受 PR 编号作为参数
   - 返回审查建议

3. **计算工具**：基础算术运算
   - 支持加、减、乘、除四则运算
   - 提供参数验证和错误处理
   - 返回计算结果

## 使用方法

### 服务端

运行服务端程序：
   ```bash
   go generator ./...
   go run ./server
   ```

### 客户端

运行客户端程序：
   ```bash
   go run ./client
   ```

## 配置说明

配置文件位于 `config/default.yaml`，包含服务器和客户端的配置信息。

## API 示例

### 1. 访问用户资源

```http
GET users://123/profile
```

响应示例：
```json
{
  "id": 10,
  "name": "Jim"
}
```

### 2. 使用计算工具

请求参数：
- operation: 运算类型 (add/subtract/multiply/divide)
- x: 第一个数字
- y: 第二个数字

示例：
```json
{
  "operation": "add",
  "x": 10,
  "y": 20
}
```

响应：
```json
{
  "result": 30
}
```

## 注意事项

1. 确保已安装 Gone 框架和相关依赖
2. 运行服务端前确保配置文件正确
3. 除零错误会在除法运算时进行检查

## 相关文档

- [Gone MCP 组件文档](../../../mcp)
- [Gone 框架文档](https://github.com/gone-io/gone)
- [Gone Viper 读取配置文档](../../../viper)
- [mcp-go](github.com/mark3labs/mcp-go)