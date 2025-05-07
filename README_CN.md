<p align="left">
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

[![license](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/gone-io/goner.jsonvalue?utm_source=godoc)](http://godoc.org/github.com/gone-io/goner)
[![Go Report Card](https://goreportcard.com/badge/github.com/gone-io/goner)](https://goreportcard.com/report/github.com/gone-io/goner)
[![codecov](https://codecov.io/gh/gone-io/goner/graph/badge.svg?token=H3CROTTDZ1)](https://codecov.io/gh/gone-io/goner)
[![Build and Test](https://github.com/gone-io/goner/actions/workflows/go.yml/badge.svg)](https://github.com/gone-io/goner/actions/workflows/go.yml)
[![Release](https://img.shields.io/github/release/gone-io/goner.svg?style=flat-square)](https://github.com/gone-io/goner/releases)

# Goner

`goner` 是 [Gone](https://github.com/gone-io/gone) 框架的官方组件库，提供了一系列可即插即用的组件，帮助开发者快速构建高质量的 Go 应用程序。

## 组件列表

### 本地配置
- [goner/viper](./viper) - 配置管理组件，基于 [spf13/viper](https://github.com/spf13/viper)

### 日志管理
- [goner/zap](./zap) - 日志组件，基于 [uber-go/zap](https://github.com/uber-go/zap)

### Web 框架
- [goner/gin](./gin) - 基于 [gin-gonic/gin](https://github.com/gin-gonic/gin) 的 Web 框架封装，提供路由管理、中间件处理、HTTP 注入等功能
- [goner/cmux](./cmux) - 基于 [soheilhy/cmux](https://github.com/soheilhy/cmux) 的多协议复用器，支持在同一端口上运行多种协议服务

### 数据库
- [goner/gorm](./gorm) - 基于 [GORM](https://gorm.io/) 的 ORM 组件，支持 MySQL、PostgreSQL、SQLite、SQL Server 和 ClickHouse 等多种数据库
- [goner/xorm](./xorm) - 基于 [XORM](https://xorm.io/) 的 ORM 组件，提供简单高效的数据库操作，支持多种数据库 

### 缓存与消息
- [goner/redis](./redis) - Redis 客户端封装，提供缓存、分布式锁等功能

### 搜索
- [goner/es](./es) - Elasticsearch 客户端封装，提供全文搜索功能

### 定时任务
- [goner/schedule](./schedule) - 定时任务组件

### 微服务
#### 配置中心
- [goner/apollo](./apollo) - 基于 [Apollo](https://www.apolloconfig.com/) 的配置中心组件，提供动态配置管理功能
- [goner/nacos](./nacos) - 基于 [Nacos](https://nacos.io/) 的配置中心组件，提供动态配置管理功能
- [goner/viper/remote](./viper/remote) - 基于多种远程配置中心（如 etcd、consul 等）的配置组件，提供统一的配置管理功能

### 注册中心
- [goner/nacos](./nacos) - 基于 [Nacos](https://nacos.io/) 的注册中心组件，提供服务注册、发现等功能
- [goner/etcd](./etcd) - 基于 [etcd](https://etcd.io/) 的注册中心组件，提供服务注册、发现等功能
- [goner/consul](./consul) - 基于 [consul](https://www.consul.io/) 的注册中心组件，提供服务注册、发现

#### RPC
- [goner/grpc](./grpc) - gRPC 客户端和服务端封装，简化微服务开发
- [goner/urllib](./urllib) - HTTP 客户端封装

### AI 组件
- [goner/openai](./openai) - OpenAI 客户端封装，提供 GPT 等 AI 能力集成
- [goner/deepseek](./deepseek) - Deepseek 客户端封装，提供国产大语言模型集成
- [goner/mcp](./mcp) - 基于 `github.com/mark3labs/mcp-go` 进行封装的工具包，它能帮助开发者快速构建 MCP (Model Context Protocol)  的服务端和客户端应用。通过使用 Gone MCP 组件，您可以轻松地将 AI 模型与您的业务系统进行集成。

### 可观测性
- [goner/otel](./otel) - OpenTelemetry 组件，提供分布式追踪、指标和日志收集功能
  - [goner/otel/tracer](./otel/tracer) - 追踪组件
    - [goner/otel/tracer/http](./otel/tracer/http) - 集成OLTP/HTTP协议的追踪Exporter
    - [goner/otel/tracer/grpc](./otel/tracer/grpc) - 集成OLTP/GRPC协议的追踪Exporter
    - [goner/otel/tracer/zipkin](./otel/tracer/zipkin) - 支持对接Zipkin的追踪Exporter

  - [goner/otel/meter](./otel/meter) - 指标组件
    - [goner/otel/meter/http](./otel/meter/http) - 集成OLTP/HTTP协议的指标Exporter
    - [goner/otel/meter/grpc](./otel/meter/grpc) - 集成OLTP/GRPC协议的指标Exporter
    - [goner/otel/meter/prometheus](./otel/meter/prometheus) - 提供Prometheus对接的Reader
		- [goner/otel/meter/prometheus/gin](./otel/meter/prometheus/gin) - 基于Gin的中间件，用于暴露Prometheus指标端点

  - [goner/otel/log](./otel/log) - 日志组件
    - [goner/otel/log/http](./otel/log/http) - 集成OLTP/HTTP协议的日志Exporter
    - [goner/otel/log/grpc](./otel/log/grpc) - 集成OLTP/GRPC协议的日志Exporter
- [goner/tracer](./tracer) - 提供程序内部traceID隐形传参



## 安装
```bash
# 安装 Gone 控制台工具
go install github.com/gone-io/gonectl@latest

# 安装 Goner 组件
# gonectl install <goner component name>
gonectl install goner/gin
```

## 快速开始
使用 Gone 控制台工具创建一个基于 Gin、XORM 和 Viper 的应用程序：

- 创建和安装依赖
```bash
gonectl create -t gin+xorm+viper goner-demo
cd goner-demo
go mod tidy
```

- 运行应用程序
```bash
# 启动数据库
docker compose up -d

# 启动应用
go run ./cmd

# 或者
# gonectl run ./cmd
```

> 示例代码，所在目录：[gin+xorm+viper](examples/gin%2Bxorm%2Bviper)


## 贡献指南

欢迎贡献代码或提出问题！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件
