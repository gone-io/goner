<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/gorm/clickhouse 组件, Gone Gorm ClickHouse 驱动

## 简介

Gone Gorm ClickHouse 驱动是 Gone Gorm 组件的 ClickHouse 数据库驱动实现。它允许您在 Gone 应用中使用 GORM 操作 ClickHouse 数据库，提供了与 Gone 框架的无缝集成。

## 特性

- 支持 ClickHouse 数据库的基本操作
- 与 Gone 框架无缝集成
- 支持连接池配置
- 提供灵活的日志配置

## 安装

```go
// 在您的应用中引入 Gone Gorm ClickHouse 驱动
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/clickhouse"
)

// 在应用初始化时加载 ClickHouse 驱动
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // 加载 Gorm 核心组件
            clickhouse.Load,  // 加载 ClickHouse 驱动
            // ...
        ).
        Run()
}
```

## 配置说明

### ClickHouse 配置

```properties
# ClickHouse 基础配置
gorm.clickhouse.driver-name=                  # 驱动名称，可选
gorm.clickhouse.dsn=tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20  # 数据源名称
gorm.clickhouse.debug=false                   # 是否启用调试模式，默认为 false
gorm.clickhouse.cluster=                      # 集群名称，可选
gorm.clickhouse.settings=                     # 其他 ClickHouse 特定设置，可选
```

## 使用示例

### 基本用法

```go
package example

import (
    "github.com/gone-io/gone/v2"
    "gorm.io/gorm"
)

type LogEntry struct {
    ID        uint      `gorm:"primaryKey"`
    Timestamp time.Time `gorm:"type:DateTime"`
    Level     string    `gorm:"type:String"`
    Message   string    `gorm:"type:String"`
}

type LogService struct {
    gone.Flag
    db *gorm.DB `gone:"*"`  // 注入 GORM 实例
}

func (s *LogService) CreateLog(level, message string) (*LogEntry, error) {
    log := &LogEntry{
        Timestamp: time.Now(),
        Level:     level,
        Message:   message,
    }
    
    if err := s.db.Create(log).Error; err != nil {
        return nil, err
    }
    
    return log, nil
}

func (s *LogService) QueryLogs(level string) ([]LogEntry, error) {
    var logs []LogEntry
    if err := s.db.Where("level = ?", level).Find(&logs).Error; err != nil {
        return nil, err
    }
    return logs, nil
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
        &LogEntry{},
        // 其他模型...
    )
}
```

## 最佳实践

1. 数据库设计
   - 合理使用 ClickHouse 的列类型
   - 根据查询模式选择适当的表引擎
   - 使用合适的分区键和排序键

2. 查询优化
   - 利用 ClickHouse 的列式存储特性
   - 避免使用 JOIN 操作
   - 合理使用物化视图

3. 写入优化
   - 使用批量插入提高性能
   - 避免频繁的小批量写入
   - 合理设置写入超时时间

4. 监控和维护
   - 监控查询性能和资源使用
   - 定期优化表结构
   - 及时清理过期数据

## 常见问题

1. **连接问题**
   
   问题：无法连接到 ClickHouse 服务器
   
   解决方案：检查网络连接、端口配置和认证信息，确保 DSN 格式正确。

2. **写入性能**
   
   问题：数据写入速度慢
   
   解决方案：使用批量插入，调整写入缓冲区大小，选择合适的表引擎。

3. **查询超时**
   
   问题：复杂查询执行超时
   
   解决方案：优化查询语句，调整超时设置，使用适当的索引和物化视图。