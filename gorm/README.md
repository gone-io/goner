# Gone Gorm 组件

## 简介

Gone Gorm 是 Gone 框架的 ORM 组件，基于 [GORM](https://gorm.io/) 实现，提供了与 Gone 框架的无缝集成。通过该组件，您可以轻松地在 Gone 应用中使用 GORM 进行数据库操作，支持 MySQL、PostgreSQL、SQLite、SQL Server 和 ClickHouse 等多种数据库。

## 特性

- 支持多种数据库：MySQL、PostgreSQL、SQLite、SQL Server 和 ClickHouse
- 与 Gone 框架无缝集成
- 支持连接池配置
- 提供灵活的日志配置
- 支持事务管理
- 支持数据库迁移

## 安装

```go
// 在您的应用中引入 Gone Gorm 组件
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    // 根据需要引入特定的数据库驱动
    _ "github.com/gone-io/goner/gorm/mysql"
    // _ "github.com/gone-io/goner/gorm/postgres"
    // _ "github.com/gone-io/goner/gorm/sqlite"
    // _ "github.com/gone-io/goner/gorm/sqlserver"
    // _ "github.com/gone-io/goner/gorm/clickhouse"
)

// 在应用初始化时加载 Gorm 组件
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // 加载 Gorm 核心组件
            gorm.mysql.Load, // 加载 MySQL 驱动
            // ...
        ).
        Run()
}
```

## 配置说明

### 基础配置

```properties
# GORM 基础配置
gorm.skip-default-transaction=false     # 是否跳过默认事务，默认为 false
gorm.full-save-associations=false       # 是否完整保存关联，默认为 false
gorm.dry-run=false                      # 是否只生成 SQL 而不执行，默认为 false
gorm.prepare-stmt=false                 # 是否使用预处理语句，默认为 false
gorm.disable-automatic-ping=false       # 是否禁用自动 ping，默认为 false
gorm.disable-foreign-key-constraint-when-migrating=false  # 迁移时是否禁用外键约束，默认为 false
gorm.ignore-relationships-when-migrating=false           # 迁移时是否忽略关系，默认为 false
gorm.disable-nested-transaction=false    # 是否禁用嵌套事务，默认为 false
gorm.allow-global-update=false           # 是否允许全局更新，默认为 false
gorm.query-fields=false                  # 是否查询所有字段，默认为 false
gorm.create-batch-size=0                 # 批量创建的大小，默认为 0
gorm.translate-error=false               # 是否翻译错误，默认为 false
gorm.propagate-unscoped=false            # 是否传播 Unscoped 到嵌套语句，默认为 false

# 连接池配置
gorm.pool.max-idle=10                    # 最大空闲连接数，默认为 10
gorm.pool.max-open=100                   # 最大打开连接数，默认为 100
gorm.pool.conn-max-lifetime=1h           # 连接最大生命周期，默认为 1 小时

# 日志配置
gorm.logger.slow-threshold=200ms         # 慢查询阈值，默认为 200ms
```

### MySQL 配置

```properties
# MySQL 基础配置
gorm.mysql.driver-name=                  # 驱动名称，可选
gorm.mysql.dsn=user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local  # 数据源名称
gorm.mysql.server-version=               # 服务器版本，可选
gorm.mysql.skip-initialize-with-version=false  # 是否跳过使用版本初始化，默认为 false
gorm.mysql.default-string-size=0         # 默认字符串大小，默认为 0
gorm.mysql.default-datetime-precision=   # 默认日期时间精度，可选
gorm.mysql.disable-with-returning=false  # 是否禁用 WITH RETURNING，默认为 false
gorm.mysql.disable-datetime-precision=false  # 是否禁用日期时间精度，默认为 false
gorm.mysql.dont-support-rename-index=false  # 是否不支持重命名索引，默认为 false
gorm.mysql.dont-support-rename-column=false  # 是否不支持重命名列，默认为 false
gorm.mysql.dont-support-for-share-clause=false  # 是否不支持 FOR SHARE 子句，默认为 false
gorm.mysql.dont-support-null-as-default-value=false  # 是否不支持 NULL 作为默认值，默认为 false
gorm.mysql.dont-support-rename-column-unique=false  # 是否不支持重命名唯一列，默认为 false
gorm.mysql.dont-support-drop-constraint=false  # 是否不支持删除约束，默认为 false
```

### PostgreSQL 配置

```properties
# PostgreSQL 基础配置
gorm.postgres.driver-name=               # 驱动名称，可选
gorm.postgres.dsn=host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai  # 数据源名称
gorm.postgres.without-quoting-check=false  # 是否不进行引号检查，默认为 false
gorm.postgres.prefer-simple-protocol=false  # 是否优先使用简单协议，默认为 false
gorm.postgres.without-returning=false    # 是否不使用 RETURNING，默认为 false
```

### SQLite 配置

```properties
# SQLite 基础配置
gorm.sqlite.driver-name=                 # 驱动名称，可选
gorm.sqlite.dsn=gorm.db                  # 数据源名称，默认为 gorm.db
```

## 使用示例

### 基本用法

```go
package example

import (
    "github.com/gone-io/gone/v2"
    "gorm.io/gorm"
)

type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:255"`
    Age  int
}

type UserService struct {
    gone.Flag
    db *gorm.DB `gone:"*"`  // 注入 GORM 实例
}

func (s *UserService) CreateUser(name string, age int) (*User, error) {
    user := &User{
        Name: name,
        Age:  age,
    }
    
    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *UserService) GetUserByID(id uint) (*User, error) {
    var user User
    if err := s.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) UpdateUser(user *User) error {
    return s.db.Save(user).Error
}

func (s *UserService) DeleteUser(id uint) error {
    return s.db.Delete(&User{}, id).Error
}
```

### 自动迁移

```go
type AppStart struct {
    gone.Flag
    db *gorm.DB `gone:"*"`
}

func (s *AppStart) AfterRevive() error {
    // 自动迁移数据库结构
    return s.db.AutoMigrate(
        &User{},
        // 其他模型...
    )
}
```

### 事务处理

```go
func (s *UserService) TransferMoney(fromID, toID uint, amount float64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 在事务中执行数据库操作
        var fromAccount, toAccount Account
        
        if err := tx.First(&fromAccount, fromID).Error; err != nil {
            return err
        }
        
        if err := tx.First(&toAccount, toID).Error; err != nil {
            return err
        }
        
        if fromAccount.Balance < amount {
            return errors.New("余额不足")
        }
        
        // 更新账户余额
        if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance - amount).Error; err != nil {
            return err
        }
        
        if err := tx.Model(&toAccount).Update("balance", toAccount.Balance + amount).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

## 最佳实践

1. 数据库连接管理
   - 合理配置连接池参数，避免连接泄漏
   - 在应用关闭时正确关闭数据库连接
   - 使用事务处理需要原子性的操作

2. 模型设计
   - 使用结构体标签定义字段属性
   - 合理设置索引和约束
   - 使用钩子函数处理模型生命周期事件

3. 查询优化
   - 使用索引优化查询性能
   - 避免 N+1 查询问题
   - 使用预加载减少查询次数

4. 日志配置
   - 设置合适的慢查询阈值
   - 在开发环境启用详细日志
   - 在生产环境适当降低日志级别

## 常见问题

1. **连接池配置**
   
   问题：应用性能下降或出现连接错误
   
   解决方案：调整连接池参数，如 `gorm.pool.max-idle` 和 `gorm.pool.max-open`，确保连接数与数据库服务器配置相匹配。

2. **慢查询问题**
   
   问题：某些查询执行时间过长
   
   解决方案：设置 `gorm.logger.slow-threshold` 以识别慢查询，然后优化这些查询，如添加索引或重写查询逻辑。

3. **事务处理**
   
   问题：事务中的错误处理不当导致数据不一致
   
   解决方案：确保在事务中正确处理错误，并在出现错误时回滚事务。

