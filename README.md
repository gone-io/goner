<p align="left">
   English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

[![license](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/gone-io/goner.jsonvalue?utm_source=godoc)](http://godoc.org/github.com/gone-io/goner)
[![Go Report Card](https://goreportcard.com/badge/github.com/gone-io/goner)](https://goreportcard.com/report/github.com/gone-io/goner)
[![codecov](https://codecov.io/gh/gone-io/goner/graph/badge.svg?token=H3CROTTDZ1)](https://codecov.io/gh/gone-io/goner)
[![Build and Test](https://github.com/gone-io/goner/actions/workflows/go.yml/badge.svg)](https://github.com/gone-io/goner/actions/workflows/go.yml)
[![Release](https://img.shields.io/github/release/gone-io/goner.svg?style=flat-square)](https://github.com/gone-io/goner/releases)

# Goner

`goner` is the official component library for the [Gone](https://github.com/gone-io/gone) framework, providing a series of plug-and-play components to help developers quickly build high-quality Go applications.

## Component List

### Web Framework
- [gin](./gin) - A web framework wrapper based on [gin-gonic/gin](https://github.com/gin-gonic/gin), providing route management, middleware processing, HTTP injection, and other features
- [cmux](./cmux) - A multi-protocol multiplexer based on [soheilhy/cmux](https://github.com/soheilhy/cmux), supporting multiple protocol services on the same port

### Database
- [gorm](./gorm) - An ORM component based on [GORM](https://gorm.io/), supporting MySQL, PostgreSQL, SQLite, SQL Server, ClickHouse, and other databases
- [xorm](./xorm) - An ORM component based on [XORM](https://xorm.io/), providing simple and efficient database operations, supporting multiple databases

### Cache and Messaging
- [redis](./redis) - Redis client wrapper, providing caching, distributed locks, and other features

### Microservices
#### Configuration Center
- [apollo](./apollo) - A configuration center component based on [Apollo](https://www.apolloconfig.com/), providing dynamic configuration management
- [nacos](./nacos) - A configuration center component based on [Nacos](https://nacos.io/), providing dynamic configuration management
- [remote](./viper/remote) - A configuration component based on various remote configuration centers (such as etcd, consul, etc.), providing unified configuration management

### Service Registry
- [nacos](./nacos) - A service registry component based on [Nacos](https://nacos.io/), providing service registration, discovery, and other features

#### RPC
- [grpc](./grpc) - gRPC client and server wrapper, simplifying microservice development
- [urllib](./urllib) - HTTP client wrapper

### AI Components
- [openai](./openai) - OpenAI client wrapper, providing GPT and other AI capabilities integration
- [deepseek](./deepseek) - Deepseek client wrapper, providing Chinese LLM integration
- [mcp](./mcp) - is a toolkit wrapped around `github.com/mark3labs/mcp-go`, helping developers quickly build MCP(Model Context Protocol)  server and client applications. By using the Gone MCP component, you can easily integrate AI models with your business systems.
- 
### Utility Components
- [viper](./viper) - Configuration management component, based on [spf13/viper](https://github.com/spf13/viper)
- [zap](./zap) - Logging component, based on [uber-go/zap](https://github.com/uber-go/zap)
- [tracer](./tracer) - Distributed tracing component
- [urllib](./urllib) - HTTP client wrapper
- [schedule](./schedule) - Scheduled task component
- [es](./es) - Elasticsearch client wrapper, providing full-text search functionality

## Installation
```bash
go get github.com/gone-io/goner
```

## Quick Start

Here's an example of creating a simple web application using the Gone framework and Goner component library:

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

// Define controller
type HelloController struct {
	gone.Flag
	gin.IRouter `gone:"*"`      // Inject router
	uR          *UserRepository `gone:"*"`
}

// Mount implements the gin.Controller interface
func (h *HelloController) Mount() gin.MountError {
	h.GET("/hello", h.hello) // Register route
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

// Define data model and repository
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
	// Load components and start the application
	gone.
		Loads(
			goner.BaseLoad,
			goneGorm.Load, // Load Gorm core components
			mysql.Load,    // Load MySQL driver
			gin.Load,      // Load Gin components
		).
		Load(&HelloController{}). // Load controller
		Load(&UserRepository{}).  // Load repository
		Serve()
}
```

- config/default.properties
```init
gorm.mysql.dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

## Configuration Guide

For detailed configuration instructions for each component, please refer to the README.md file in each component directory.

## Contribution Guidelines

Contributions of code or issues are welcome! Please follow these steps:

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details