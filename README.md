<p align="left">
    <a href="README_CN.md">中文</a>&nbsp;|&nbsp;English
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

### Local Configuration
- [goner/viper](./viper) - Configuration management component based on [spf13/viper](https://github.com/spf13/viper)

### Logging
- [goner/zap](./zap) - Logging component based on [uber-go/zap](https://github.com/uber-go/zap)

### Web Framework
- [goner/gin](./gin) - Web framework wrapper based on [gin-gonic/gin](https://github.com/gin-gonic/gin), providing route management, middleware processing, HTTP injection and other functions
- [goner/cmux](./cmux) - Multi-protocol multiplexer based on [soheilhy/cmux](https://github.com/soheilhy/cmux), supporting multiple protocol services on the same port

### Database
- [goner/gorm](./gorm) - ORM component based on [GORM](https://gorm.io/), supporting MySQL, PostgreSQL, SQLite, SQL Server and ClickHouse
- [goner/xorm](./xorm) - ORM component based on [XORM](https://xorm.io/), providing simple and efficient database operations, supporting multiple databases

### Cache & Messaging
- [goner/redis](./redis) - Redis client wrapper, providing caching, distributed lock and other functions

### Search
- [goner/es](./es) - Elasticsearch client wrapper, providing full-text search functionality

### Scheduling
- [goner/schedule](./schedule) - Scheduled task component

### Microservices
#### Configuration Center
- [goner/apollo](./apollo) - Configuration center component based on [Apollo](https://www.apolloconfig.com/), providing dynamic configuration management
- [goner/nacos](./nacos) - Configuration center component based on [Nacos](https://nacos.io/), providing dynamic configuration management
- [goner/viper/remote](./viper/remote) - Configuration component based on various remote configuration centers (such as etcd, consul), providing unified configuration management

### Service Registry
- [goner/nacos](./nacos) - Service registry component based on [Nacos](https://nacos.io/), providing service registration and discovery
- [goner/etcd](./etcd) - Service registry component based on [etcd](https://etcd.io/), providing service registration and discovery
- [goner/consul](./consul) - Service registry component based on [consul](https://www.consul.io/), providing service registration and discovery

#### RPC
- [goner/grpc](./grpc) - gRPC client and server wrapper, simplifying microservice development
- [goner/urllib](./urllib) - HTTP client wrapper

### AI Components
- [goner/openai](./openai) - OpenAI client wrapper, providing GPT and other AI capabilities integration
- [goner/deepseek](./deepseek) - Deepseek client wrapper, providing domestic large language model integration
- [goner/mcp](./mcp) - Toolkit wrapper based on `github.com/mark3labs/mcp-go`, helping developers quickly build MCP (Model Context Protocol) server and client applications

### Observability
- [goner/otel](./otel) - OpenTelemetry component, providing distributed tracing, metrics and log collection
  - [goner/otel/tracer](./otel/tracer) - Tracing component
    - [goner/otel/tracer/http](./otel/tracer/http) - Tracing Exporter integrated with OLTP/HTTP protocol
    - [goner/otel/tracer/grpc](./otel/tracer/grpc) - Tracing Exporter integrated with OLTP/GRPC protocol
    - [goner/otel/tracer/zipkin](./otel/tracer/zipkin) - Tracing Exporter supporting Zipkin

  - [goner/otel/meter](./otel/meter) - Metrics component
    - [goner/otel/meter/http](./otel/meter/http) - Metrics Exporter integrated with OLTP/HTTP protocol
    - [goner/otel/meter/grpc](./otel/meter/grpc) - Metrics Exporter integrated with OLTP/GRPC protocol
    - [goner/otel/meter/prometheus](./otel/meter/prometheus) - Reader providing Prometheus integration
    - [goner/otel/meter/prometheus/gin](./otel/meter/prometheus/gin) - Gin middleware for exposing Prometheus metrics endpoint

  - [goner/otel/log](./otel/log) - Log component
    - [goner/otel/log/http](./otel/log/http) - Log Exporter integrated with OLTP/HTTP protocol
    - [goner/otel/log/grpc](./otel/log/grpc) - Log Exporter integrated with OLTP/GRPC protocol
- [goner/tracer](./tracer) - Providing implicit parameter passing of traceID within programs

## Installation
```bash
# Install Gone CLI tool
go install github.com/gone-io/gonectl@latest

# Install Goner component
# gonectl install <goner component name>
gonectl install goner/gin
```

## Quick Start
Create an application based on Gin, XORM and Viper using Gone CLI tool:

- Create and install dependencies
```bash
gonectl create -t gin+xorm+viper goner-demo
cd goner-demo
go mod tidy
```

- Run application
```bash
# Start database
docker compose up -d

# Start application
go run ./cmd

# Or
# gonectl run ./cmd
```

> Example code location: [gin+xorm+viper](examples/gin%2Bxorm%2Bviper)


## Contribution Guide

Welcome contributions! Please follow these steps:

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file