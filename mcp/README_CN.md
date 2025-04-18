# Gone MCP 组件使用指南


## 特性：
- 服务端
  - 支持多实例
  - 支持 `Goner Define` 和 `Functional Define` 定义 MCP 的 Tool、Prompt、Resource
  - 支持配置文件
  - 支持定义Hooks注入
  - 支持定义ContextFunc注入
  - 支持Stdio和SEE
- 客户端
  - 支持多实例，根据不同的`gone`标签配置获取不同的实例
  - 支持从配置中读取参数
  - 支持Stdio和SEE
  - SSE支持设置header
  - 支持自定义`transport.Interface`