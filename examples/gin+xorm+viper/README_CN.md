<p align="left">
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone Web MySQL 模板项目

这是一个基于 Gone 框架的 Web 应用模板项目，集成了 MySQL 数据库支持，可以作为开发 Web 应用的起点。

## 项目特点

- 基于 Gone 框架构建，采用依赖注入设计模式
- 集成 MySQL 数据库支持
- 提供完整的项目结构和代码组织方式
- 支持 Docker 部署
- 包含用户管理示例代码

## 技术栈

- [Gone](https://github.com/gone-io/gone) - 基于依赖注入的 Go 应用框架
- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [Xorm](https://xorm.io/) - ORM 框架
- MySQL - 数据库
- Docker - 容器化部署

## 目录结构

```
.
├── cmd                 # 应用入口
│   └── server          # 服务器入口
├── config              # 配置文件
├── internal            # 内部代码
│   ├── controller      # 控制器
│   ├── interface       # 接口定义
│   ├── module          # 模块实现
│   ├── pkg             # 工具包
│   └── router          # 路由定义
├── scripts             # 脚本文件
│   └── mysql           # MySQL 初始化脚本
└── tests               # 测试文件
    └── api             # API 测试
```

## 快速开始

### 前置条件

- Go 1.24 或更高版本
- MySQL 8.0 或更高版本
- Docker 和 Docker Compose (可选，用于容器化部署)

### 本地开发

1. 安装依赖

```bash
go mod download
```

3. 配置数据库

编辑 `config/default.properties` 文件，设置数据库连接信息：

```properties
db.host=localhost
db.port=3306
db.name=demo
db.username=root
db.password=123456
```

4. 运行项目

```bash
make run
```

服务将在 http://localhost:8080 启动

### Docker 部署

1. 构建 Docker 镜像

```bash
make build-docker
```

2. 启动服务

```bash
docker compose up -d
```

服务将在 http://localhost:8080 启动，MySQL 数据库将在 localhost:3306 可用

## 配置说明

项目配置文件位于 `config/default.properties`，主要配置项包括：

- `server.port` - 服务端口，默认为 8080
- `server.mode` - Gin 服务模式，可选值 debug, test, release，默认为 release
- `server.health-check` - 健康检查路径，默认为 /api/health-check
- `db.host` - 数据库主机
- `db.port` - 数据库端口
- `db.name` - 数据库名称
- `db.username` - 数据库用户名
- `db.password` - 数据库密码

## 开发指南

### 添加新接口

1. 在 `internal/interface` 中定义接口
2. 在 `internal/module` 中实现接口
3. 在 `internal/controller` 中创建控制器
4. 在 `internal/router` 中注册路由

### 生成 Mock

项目使用 mockgen 生成 mock 代码：

```bash
go generate ./...
```