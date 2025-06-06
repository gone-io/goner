<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

[![license](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/gone-io/goner.jsonvalue?utm_source=godoc)](http://godoc.org/github.com/gone-io/goner)
[![Go Report Card](https://goreportcard.com/badge/github.com/gone-io/goner)](https://goreportcard.com/report/github.com/gone-io/goner)
[![codecov](https://codecov.io/gh/gone-io/goner/graph/badge.svg?token=H3CROTTDZ1)](https://codecov.io/gh/gone-io/goner)
[![Build and Test](https://github.com/gone-io/goner/actions/workflows/go.yml/badge.svg)](https://github.com/gone-io/goner/actions/workflows/go.yml)
[![Release](https://img.shields.io/github/release/gone-io/goner.svg?style=flat-square)](https://github.com/gone-io/goner/releases)

# Goner

Goner is the official component library for the [Gone](https://github.com/gone-io/gone) framework, providing a series of plug-and-play components to help developers quickly build high-quality Go applications.

## Component List

- Configuration Management
    - Local Configuration
        - [goner/viper](./viper) - Configuration management component based on [spf13/viper](https://github.com/spf13/viper)
    - Remote Configuration/Configuration Center [Microservices]
        - [goner/apollo](./apollo) - Configuration center component based on [Apollo](https://www.apolloconfig.com/), providing dynamic configuration management
        - [goner/nacos](./nacos) - Configuration center component based on [Nacos](https://nacos.io/), providing dynamic configuration management
        - [goner/viper/remote](./viper/remote) - Configuration component based on various remote configuration centers (such as etcd, consul, etc.), providing unified configuration management
- Log Management
    - [goner/zap](./zap) - Logging component based on [uber-go/zap](https://github.com/uber-go/zap)

- Service Registry [Microservices]
    - [goner/nacos](./nacos) - Service registry component based on [Nacos](https://nacos.io/), providing service registration, discovery, and other features
    - [goner/etcd](./etcd) - Service registry component based on [etcd](https://etcd.io/), providing service registration, discovery, and other features
    - [goner/consul](./consul) - Service registry component based on [consul](https://www.consul.io/), providing service registration and discovery

- Message Queue [Microservices] [Event Storming]
    - [goner/mq/kafka](./mq/kafka) - Provides Kafka integration
    - [goner/mq/rocket](./mq/rocket) - Provides RocketMQ integration
    - [goner/mq/rabbitmq](./mq/rabbitmq) - Provides RabbitMQ integration
    - [goner/mq/mqtt](./mq/mqtt) - Provides MQTT integration

- Services and Remote Calls [Microservices]
    - [goner/grpc](./grpc) - gRPC client and server wrapper, simplifying microservice development
    - [goner/urllib](./urllib) - HTTP client wrapper
    - [goner/gin](./gin) - Web framework wrapper based on [gin-gonic/gin](https://github.com/gin-gonic/gin), providing route management, middleware handling, HTTP injection, and other features
    - [goner/cmux](./cmux) - Multi-protocol multiplexer based on [soheilhy/cmux](https://github.com/soheilhy/cmux), supporting multiple protocol services on the same port

- Database
    - Relational Database
        - [goner/xorm](./xorm) - ORM component based on [XORM](https://xorm.io/), providing simple and efficient database operations, supporting multiple databases
            - [goner/xorm/mysql](./xorm/mysql) - MySQL driver wrapper for Xorm, providing database operation functionality
            - [goner/xorm/postgres](./xorm/postgres) - PostgreSQL driver wrapper for Xorm, providing database operation functionality
            - [goner/xorm/sqlite](./xorm/sqlite) - SQLite driver wrapper for Xorm, providing database operation functionality
            - [goner/xorm/mssql](./xorm/mssql) - MSSQL driver wrapper for Xorm, providing database operation functionality
        - [goner/gorm](./gorm) - ORM component based on [GORM](https://gorm.io/), supporting MySQL, PostgreSQL, SQLite, SQL Server, and ClickHouse
            - [goner/gorm/mysql](./gorm/mysql) - MySQL driver wrapper for Gorm, providing database operation functionality
            - [goner/gorm/postgres](./gorm/postgres) - PostgreSQL driver wrapper for Gorm, providing database operation functionality
            - [goner/gorm/sqlite](./gorm/sqlite) - SQLite driver wrapper for Gorm, providing database operation functionality
            - [goner/gorm/clickhouse](./gorm/clickhouse) - ClickHouse driver wrapper for Gorm, providing database operation functionality
            - [goner/gorm/sqlserver](./gorm/sqlserver) - SqlServer driver wrapper for Gorm, providing database operation functionality
    - NoSQL
        - [goner/redis](./redis) - Redis client wrapper, providing caching, distributed locking, and other features
        - [goner/es](./es) - Elasticsearch client wrapper, providing full-text search functionality

- Observability [Microservices]
    - [goner/otel](./otel) - OpenTelemetry component, providing distributed tracing, metrics, and log collection
        - [goner/otel/tracer](./otel/tracer) - Tracing component
            - [goner/otel/tracer/http](./otel/tracer/http) - OLTP/HTTP protocol tracing Exporter integration
            - [goner/otel/tracer/grpc](./otel/tracer/grpc) - OLTP/GRPC protocol tracing Exporter integration
            - [goner/otel/tracer/zipkin](./otel/tracer/zipkin) - Zipkin tracing Exporter support

        - [goner/otel/meter](./otel/meter) - Metrics component
            - [goner/otel/meter/http](./otel/meter/http) - OLTP/HTTP protocol metrics Exporter integration
            - [goner/otel/meter/grpc](./otel/meter/grpc) - OLTP/GRPC protocol metrics Exporter integration
            - [goner/otel/meter/prometheus](./otel/meter/prometheus) - Prometheus Reader integration
                - [goner/otel/meter/prometheus/gin](./otel/meter/prometheus/gin) - Gin middleware for exposing Prometheus metrics endpoints

        - [goner/otel/log](./otel/log) - Logging component
            - [goner/otel/log/http](./otel/log/http) - OLTP/HTTP protocol logging Exporter integration
            - [goner/otel/log/grpc](./otel/log/grpc) - OLTP/GRPC protocol logging Exporter integration
    - [goner/tracer](./tracer) - Provides internal traceID implicit parameter passing

- Scheduled Tasks
    - [goner/schedule](./schedule) - Scheduled task component

- AI Components [Large Language Models] [Artificial Intelligence]
    - [goner/openai](./openai) - OpenAI client wrapper, providing GPT and other AI capabilities integration
    - [goner/deepseek](./deepseek) - Deepseek client wrapper, providing domestic large language model integration
    - [goner/mcp](./mcp) - A toolkit wrapped based on `github.com/mark3labs/mcp-go` that helps developers quickly build MCP (Model Context Protocol) server and client applications. Using the Gone MCP component, you can easily integrate AI models with your business systems.

## Installation

```bash
# Install Gone CLI tool
go install github.com/gone-io/gonectl@latest

# Install Goner components
# gonectl install <goner component name>
gonectl install goner/gin
```

## Quick Start

Create an application based on Gin, XORM, and Viper using the Gone CLI tool:

- Create and install dependencies

```bash
gonectl create -t gin+xorm+viper goner-demo
cd goner-demo
go mod tidy
```

- Run the application

```bash
# Start the database
docker compose up -d

# Start the application
go run ./cmd

# Or
# gonectl run ./cmd
```

> Sample code location: [gin+xorm+viper](examples/gin%2Bxorm%2Bviper)

## Contributing Guidelines

Contributions of code or issues are welcome! Please follow these steps:

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details
