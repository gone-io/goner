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

### Web 框架
- [gin](./gin) - 基于 [gin-gonic/gin](https://github.com/gin-gonic/gin) 的 Web 框架封装，提供路由管理、中间件处理、HTTP 注入等功能
- [cmux](./cmux) - 基于 [soheilhy/cmux](https://github.com/soheilhy/cmux) 的多协议复用器，支持在同一端口上运行多种协议服务

### 数据库
- [gorm](./gorm) - 基于 [GORM](https://gorm.io/) 的 ORM 组件，支持 MySQL、PostgreSQL、SQLite、SQL Server 和 ClickHouse 等多种数据库
- [xorm](./xorm) - 基于 [XORM](https://xorm.io/) 的 ORM 组件，提供简单高效的数据库操作，支持多种数据库 

### 缓存与消息
- [redis](./redis) - Redis 客户端封装，提供缓存、分布式锁等功能

### 微服务
#### 配置中心
- [apollo](./apollo) - 基于 [Apollo](https://www.apolloconfig.com/) 的配置中心组件，提供动态配置管理功能
- [nacos](./nacos) - 基于 [Nacos](https://nacos.io/) 的配置中心组件，提供动态配置管理功能
- [remote](./viper/remote) - 基于多种远程配置中心（如 etcd、consul 等）的配置组件，提供统一的配置管理功能

### 注册中心
- [nacos](./nacos) - 基于 [Nacos](https://nacos.io/) 的注册中心组件，提供服务注册、发现等功能
- [etcd](./etcd)] - 基于 [etcd](https://etcd.io/) 的注册中心组件，提供服务注册、发现等功能
- [consul](./consul) - 基于 [consul](https://www.consul.io/) 的注册中心组件，提供服务注册、发现

#### RPC
- [grpc](./grpc) - gRPC 客户端和服务端封装，简化微服务开发
- [urllib](./urllib) - HTTP 客户端封装

### AI 组件
- [openai](./openai) - OpenAI 客户端封装，提供 GPT 等 AI 能力集成
- [deepseek](./deepseek) - Deepseek 客户端封装，提供国产大语言模型集成
- [mcp](./mcp) - 基于 `github.com/mark3labs/mcp-go` 进行封装的工具包，它能帮助开发者快速构建 MCP (Model Context Protocol)  的服务端和客户端应用。通过使用 Gone MCP 组件，您可以轻松地将 AI 模型与您的业务系统进行集成。

### 工具组件
- [viper](./viper) - 配置管理组件，基于 [spf13/viper](https://github.com/spf13/viper)
- [zap](./zap) - 日志组件，基于 [uber-go/zap](https://github.com/uber-go/zap)
- [tracer](./tracer) - 分布式追踪组件
- [urllib](./urllib) - HTTP 客户端封装
- [schedule](./schedule) - 定时任务组件
- [es](./es) - Elasticsearch 客户端封装，提供全文搜索功能

## 安装
```bash
go get github.com/gone-io/goner
```

## 快速开始

以下是一个使用 Gone 框架和 Goner 组件库创建简单 Web 应用的示例：


- main.go
```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
	goneGorm "github.com/gone-io/goner/gorm"
	"github.com/gone-io/goner/gorm/mysql"
	"gorm.io/gorm"
)

// 定义控制器
type HelloController struct {
	gone.Flag
	gin.IRouter `gone:"*"`      // 注入路由器
	uR          *UserRepository `gone:"*"`
}

// Mount 实现 gin.Controller 接口
func (h *HelloController) Mount() gin.MountError {
	h.GET("/hello", h.hello) // 注册路由
	h.GET("/user/:id", h.getUser)
	return nil
}

func (h *HelloController) hello() (string, error) {
	return "Hello, Gone!", nil
}
func (h *HelloController) getUser(in struct {
	id uint `param:"id"`
}) (*User, error) {

	user, err := h.uR.GetByID(in.id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 定义数据模型和仓库
type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type UserRepository struct {
	gone.Flag
	*gorm.DB `gone:"*"`
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	var user User
	err := r.First(&user, id).Error
	return &user, err
}

func main() {
	// 加载组件并启动应用
	gone.
		Loads(
			goner.BaseLoad,
			goneGorm.Load, // 加载 Gorm 核心组件
			mysql.Load,    // 加载 MySQL 驱动
			gin.Load,      // 加载 Gin 组件
		).
		Load(&HelloController{}). // 加载控制器
		Load(&UserRepository{}).  // 加载仓库
		Serve()
}
```

- config/default.properties
```init
gorm.mysql.dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

## 配置说明

各组件的详细配置说明请参考各组件目录下的 README.md 文件。

## 贡献指南

欢迎贡献代码或提出问题！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件
