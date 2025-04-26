<p align="left">
   English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Web MySQL Template Project

This is a web application template project based on the Gone framework, integrated with MySQL database support, which can serve as a starting point for developing web applications.

## Project Features

- Built on the Gone framework using dependency injection design pattern
- Integrated MySQL database support
- Provides complete project structure and code organization
- Supports Docker deployment
- Includes user management example code

## Tech Stack

- [Gone](https://github.com/gone-io/gone) - Go application framework based on dependency injection
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Xorm](https://xorm.io/) - ORM framework
- MySQL - Database
- Docker - Containerized deployment

## Directory Structure

```
.
├── cmd                 # Application entry
│   └── server          # Server entry
├── config              # Configuration files
├── internal            # Internal code
│   ├── controller      # Controllers
│   ├── interface       # Interface definitions
│   ├── module          # Module implementations
│   ├── pkg             # Utility packages
│   └── router          # Route definitions
├── scripts             # Script files
│   └── mysql           # MySQL initialization scripts
└── tests               # Test files
    └── api             # API tests
```

## Quick Start

### Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher
- Docker and Docker Compose (optional, for containerized deployment)

### Local Development

1. Install dependencies

```bash
go mod download
```

3. Configure database

Edit `config/default.properties` file, set database connection information:

```properties
db.host=localhost
db.port=3306
db.name=demo
db.username=root
db.password=123456
```

4. Run the project

```bash
make run
```

The service will start at http://localhost:8080

### Docker Deployment

1. Build Docker image

```bash
make build-docker
```

2. Start service

```bash
docker compose up -d
```

The service will start at http://localhost:8080, and MySQL database will be available at localhost:3306

## Configuration Guide

Project configuration file is located at `config/default.properties`, main configuration items include:

- `server.port` - Service port, default is 8080
- `server.mode` - Gin service mode, options are debug, test, release, default is release
- `server.health-check` - Health check path, default is /api/health-check
- `db.host` - Database host
- `db.port` - Database port
- `db.name` - Database name
- `db.username` - Database username
- `db.password` - Database password

## Development Guide

### Adding New APIs

1. Define interface in `internal/interface`
2. Implement interface in `internal/module`
3. Create controller in `internal/controller`
4. Register route in `internal/router`

### Generating Mocks

The project uses mockgen to generate mock code:

```bash
go generate ./...
```