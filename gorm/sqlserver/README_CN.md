# Gone Gorm SQL Server 驱动

## 简介

Gone Gorm SQL Server 驱动是 Gone Gorm 组件的 Microsoft SQL Server 数据库驱动实现。它允许您在 Gone 应用中使用 GORM 操作 SQL Server 数据库，提供了与 Gone 框架的无缝集成。

## 特性

- 支持 SQL Server 数据库的基本操作
- 与 Gone 框架无缝集成
- 支持连接池配置
- 提供灵活的日志配置
- 支持事务管理
- 支持数据库迁移

## 安装

```go
// 在您的应用中引入 Gone Gorm SQL Server 驱动
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/sqlserver"
)

// 在应用初始化时加载 SQL Server 驱动
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // 加载 Gorm 核心组件
            sqlserver.Load,   // 加载 SQL Server 驱动
            // ...
        ).
        Run()
}
```

## 配置说明

### SQL Server 配置

```properties
# SQL Server 基础配置
gorm.sqlserver.driver-name=              # 驱动名称，可选
gorm.sqlserver.dsn=sqlserver://gorm:gorm@localhost:9930?database=gorm  # 数据源名称
gorm.sqlserver.schema=                   # 数据库架构，默认为 dbo
gorm.sqlserver.encrypt=disable           # 加密连接，可选值：disable, true, false
gorm.sqlserver.trust-server-certificate=false  # 是否信任服务器证书，默认为 false
gorm.sqlserver.app-name=                 # 应用名称，可选
```

## 使用示例

### 基本用法

```go
package example

import (
    "github.com/gone-io/gone/v2"
    "gorm.io/gorm"
)

type Employee struct {
    ID        uint      `gorm:"primaryKey"`
    FirstName string    `gorm:"size:50"`
    LastName  string    `gorm:"size:50"`
    Email     string    `gorm:"size:255;uniqueIndex"`
    HireDate  time.Time `gorm:"type:datetime2"`
}

type EmployeeService struct {
    gone.Flag
    db *gorm.DB `gone:"*"` // 注入 GORM 实例
}

func (s *EmployeeService) CreateEmployee(employee *Employee) error {
    return s.db.Create(employee).Error
}

func (s *EmployeeService) GetEmployeeByID(id uint) (*Employee, error) {
    var employee Employee
    if err := s.db.First(&employee, id).Error; err != nil {
        return nil, err
    }
    return &employee, nil
}

func (s *EmployeeService) UpdateEmployee(employee *Employee) error {
    return s.db.Save(employee).Error
}

func (s *EmployeeService) DeleteEmployee(id uint) error {
    return s.db.Delete(&Employee{}, id).Error
}
```

### 使用存储过程

```go
func (s *EmployeeService) GetEmployeesByDepartment(deptID int) ([]Employee, error) {
    var employees []Employee
    err := s.db.Raw("EXEC GetEmployeesByDepartment ?", deptID).Scan(&employees).Error
    return employees, err
}
```

## 最佳实践

1. 数据库设计
   - 使用适当的数据类型
   - 合理设计索引
   - 遵循 SQL Server 命名规范
   - 使用适当的架构管理

2. 连接管理
   - 配置合适的连接池参数
   - 使用连接加密保护数据
   - 监控连接状态

3. 性能优化
   - 使用适当的索引
   - 优化查询语句
   - 合理使用存储过程
   - 避免过度使用 ORM 功能

4. 安全性
   - 使用最小权限原则
   - 启用连接加密
   - 定期更新密码
   - 避免 SQL 注入

## 常见问题

1. **连接问题**
   
   问题：无法连接到 SQL Server
   
   解决方案：检查网络连接、认证信息和防火墙设置，确保 DSN 格式正确。

2. **性能问题**
   
   问题：查询执行缓慢
   
   解决方案：检查索引使用情况，优化查询语句，使用查询计划分析器。

3. **字符集问题**
   
   问题：特殊字符显示乱码
   
   解决方案：确保数据库和连接字符集配置正确，使用 NVARCHAR 类型存储 Unicode 字符。