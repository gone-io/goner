# MongoDB Basic Example

这是一个使用 Gone 框架和 MongoDB 的基础示例。

## 快速开始

### 1. 启动 MongoDB 服务

使用 Docker Compose 启动 MongoDB 和 Mongo Express（Web 管理界面）：

```bash
docker-compose up -d
```

这将启动以下服务：
- **MongoDB**: 运行在端口 27017
- **Mongo Express**: Web 管理界面，运行在端口 8081

### 2. 访问 Mongo Express

打开浏览器访问 [http://localhost:8081](http://localhost:8081) 来管理 MongoDB 数据库。

### 3. 运行示例应用

```bash
go run .
```

## 配置说明

### MongoDB 连接信息

- **主机**: localhost
- **端口**: 27017
- **数据库**: myapp
- **管理员用户**: admin
- **管理员密码**: password
- **应用用户**: appuser
- **应用密码**: apppassword

### 配置文件

应用的配置在 `config/default.yaml` 中。默认配置已经设置为连接本地 MongoDB：

```yaml
mongo:
  uri: "mongodb://localhost:27017"
  database: "myapp"
```

如果需要使用认证，可以修改 URI 为：
```yaml
mongo:
  uri: "mongodb://appuser:apppassword@localhost:27017/myapp"
```

## 示例功能

这个示例演示了以下 MongoDB 操作：

1. **创建用户** - 插入新的用户文档
2. **查询用户** - 根据邮箱查找用户
3. **更新用户** - 修改用户信息
4. **列出用户** - 分页查询所有用户
5. **删除用户** - 删除指定用户

## 数据结构

用户文档结构：

```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Age       int                `bson:"age" json:"age"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

## 停止服务

```bash
docker-compose down
```

如果需要删除数据卷：

```bash
docker-compose down -v
```