<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/xorm 组件

**goner/xorm** 组件 是一个为 Gone 框架提供的强大数据库集成组件，基于 XORM 提供增强功能。

## 特性

- 多数据库支持（MySQL、SQLite3、PostgreSQL、Oracle、MSSQL）
- 自动事务管理
- 事务传播
- 命名参数支持
- 增强的 SQL 查询功能
- 主从集群支持
- 连接池管理

## 安装

```bash
go get github.com/gone-io/goner/xorm
```

## 导入所需的数据库驱动：

```go
import (
    github.com/gone-io/goner/xorm

    // MySQL driver
    _ "github.com/go-sql-driver/mysql"

    // SQLite3 driver
    // _ "github.com/mattn/go-sqlite3"

    // PostgreSQL driver
    // _ "github.com/lib/pq"

    // Oracle driver
    // _ "github.com/mattn/go-oci8"

    // MSSQL driver
    // _ "github.com/denisenkom/go-mssqldb"
)
```

## 配置

### 单数据库模式

```ini
database.driver-name=mysql
database.dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

### 集群模式

```ini
database.cluster.enable=true

# 主数据库配置
database.cluster.master.driver-name=mysql
database.cluster.master.dsn=root:123456@tcp(master-db-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local

# 从数据库配置
database.cluster.slaves[0].driver-name=mysql
database.cluster.slaves[0].dsn=root:123456@tcp(slave-db-0-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

### 配置参数

| 参数 | 是否必需 | 默认值 | 描述 |
|-----------|----------|---------|-------------|
| database.cluster.enable | 否 | false | 启用集群模式 |
| database.driver-name | 否* | - | 数据库驱动名称 |
| database.dsn | 否* | - | 数据库连接字符串 |
| database.max-idle-count | 否 | 5 | 最大空闲连接数 |
| database.max-open | 否 | 20 | 最大打开连接数 |
| database.max-lifetime | 否 | 10m | 连接最大生命周期 |
| database.show-sql | 否 | true | 显示 SQL 日志 |

*在非集群模式下必需

## 加载和启动
```go
package main
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/gone/goner/xorm"
    _ "github.com/go-sql-driver/mysql"
)
func main() {
    gone.
        NewApp(
            xorm.Load,
        ).
        Run(func(db xorm.Engine){
            // 你的应用逻辑
            db.Get(/*...*/)
        })
}
```

## 使用方法

### 基本用法

```go
type dbUser struct {
    gone.Flag
    db xorm.Engine `gone:"*"` // 注入数据库引擎
}

type Book struct {
    Id    int64
    Title string
}

func (d *dbUser) GetBookById(id int64) (book *Book, err error) {
    book = new(Book)
    has, err := d.db.Where("id=?", id).Get(book)
    if err != nil {
        return nil, gone.ToError(err)
    }
    if !has {
        return nil, gone.NewParameterError("book not found", 404)
    }
    return book, nil
}
```

### 自动事务

```go
func (d *db) updateUser(user *entity.User) error {
    return d.Transaction(func(session xorm.Interface) error {
        _, err := session.ID(user.Id).Update(user)
        return gone.ToError(err)
    })
}
```

### 命名参数

```go
sql, args := xorm.MustNamed(`
    update user
    set status = :status,
        avatar = :avatar
    where id = :id`,
    map[string]any{
        "id":     1,
        "status": 1,
        "avatar": "https://example.com/avatar.jpg",
    },
)
```

## 最佳实践

1. 连接池配置
   - 根据负载调整 `max-idle-count`
   - 设置适当的 `max-open` 以防止数据库过载
   - 配置 `max-lifetime` 以管理连接生命周期

2. 集群模式优化
   - 有效使用读写分离
   - 分别配置主从数据库的连接池
   - 考虑主从复制延迟对事务的影响

3. SQL 优化
   - 在开发环境中使用 `show-sql` 进行 SQL 调试
   - 实施适当的索引
   - 避免大事务
   - 使用批量操作代替循环

## 许可证

MIT 许可证