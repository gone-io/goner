# Gone Gorm MySQL 驱动

## 简介

Gone Gorm MySQL 驱动是 Gone Gorm 组件的 MySQL 数据库驱动实现。它允许您在 Gone 应用中使用 GORM 操作 MySQL 数据库，提供了与 Gone 框架的无缝集成。

## 特性

- 支持 MySQL 数据库的基本操作
- 与 Gone 框架无缝集成
- 支持连接池配置
- 提供灵活的日志配置
- 支持事务管理
- 支持数据库迁移

## 安装

```go
// 在您的应用中引入 Gone Gorm MySQL 驱动
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/mysql"
)

// 在应用初始化时加载 MySQL 驱动
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // 加载 Gorm 核心组件
            mysql.Load,       // 加载 MySQL 驱动
            // ...
        ).
        Run()
}
```

## 配置说明

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
    db *gorm.DB `gone:"*"` // 注入 GORM 实例
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
```

### 事务处理

```go
func (s *UserService) TransferMoney(fromID, toID uint, amount float64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
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

1. 数据库设计
   - 合理设计表结构和字段类型
   - 使用适当的索引提升查询性能
   - 遵循数据库设计规范

2. 连接管理
   - 合理配置连接池参数
   - 及时关闭不需要的连接
   - 使用连接池监控工具

3. 查询优化
   - 使用适当的索引
   - 避免全表扫描
   - 优化 JOIN 查询
   - 使用预处理语句

4. 事务处理
   - 合理使用事务隔离级别
   - 避免长事务
   - 正确处理事务回滚

## 常见问题

1. **连接问题**
   
   问题：无法连接到 MySQL 服务器
   
   解决方案：检查网络连接、端口配置和认证信息，确保 DSN 格式正确。

2. **字符集问题**
   
   问题：中文或特殊字符显示乱码
   
   解决方案：确保使用 utf8mb4 字符集，并在 DSN 中正确配置 charset 参数。

3. **性能问题**
   
   问题：查询执行缓慢
   
   解决方案：检查索引使用情况，优化查询语句，调整数据库配置参数。