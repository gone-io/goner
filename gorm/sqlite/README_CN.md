<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/gorm/sqlite 组件，Gone Gorm SQLite 驱动

## 简介

Gone Gorm SQLite 驱动是 Gone Gorm 组件的 SQLite 数据库驱动实现。它允许您在 Gone 应用中使用 GORM 操作 SQLite 数据库，提供了与 Gone 框架的无缝集成。SQLite 是一个轻量级的、零配置的、自包含的数据库引擎，特别适合嵌入式应用和开发环境。

## 特性

- 支持 SQLite 数据库的基本操作
- 与 Gone 框架无缝集成
- 零配置，易于使用
- 支持事务管理
- 支持数据库迁移
- 适合开发和测试环境
- 支持内存数据库

## 安装

```go
// 在您的应用中引入 Gone Gorm SQLite 驱动
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/sqlite"
)

// 在应用初始化时加载 SQLite 驱动
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // 加载 Gorm 核心组件
            sqlite.Load,      // 加载 SQLite 驱动
            // ...
        ).
        Run()
}
```

## 配置说明

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

type Note struct {
    ID      uint   `gorm:"primaryKey"`
    Title   string `gorm:"size:255"`
    Content string `gorm:"type:text"`
}

type NoteService struct {
    gone.Flag
    db *gorm.DB `gone:"*"` // 注入 GORM 实例
}

func (s *NoteService) CreateNote(title, content string) (*Note, error) {
    note := &Note{
        Title:   title,
        Content: content,
    }
    
    if err := s.db.Create(note).Error; err != nil {
        return nil, err
    }
    
    return note, nil
}

func (s *NoteService) GetNoteByID(id uint) (*Note, error) {
    var note Note
    if err := s.db.First(&note, id).Error; err != nil {
        return nil, err
    }
    return &note, nil
}

func (s *NoteService) UpdateNote(note *Note) error {
    return s.db.Save(note).Error
}

func (s *NoteService) DeleteNote(id uint) error {
    return s.db.Delete(&Note{}, id).Error
}
```

### 使用内存数据库

```go
// 配置内存数据库
gorm.sqlite.dsn=:memory:

// 或者在代码中直接使用
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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
        &Note{},
        // 其他模型...
    )
}
```

## 最佳实践

1. 数据库设计
   - 合理设计表结构
   - 适当使用索引
   - 避免过度复杂的关系

2. 性能优化
   - 使用事务处理批量操作
   - 定期进行 VACUUM 操作
   - 合理设置 journal_mode

3. 并发处理
   - 注意 SQLite 的并发限制
   - 合理使用事务隔离级别
   - 避免长时间锁定

4. 备份和维护
   - 定期备份数据库文件
   - 监控数据库大小
   - 及时清理不需要的数据

## 常见问题

1. **并发访问问题**
   
   问题：多个连接同时访问数据库导致锁定
   
   解决方案：使用适当的锁定策略，避免长事务，考虑使用 WAL 模式。

2. **性能问题**
   
   问题：数据库操作速度变慢
   
   解决方案：定期执行 VACUUM，优化索引，使用适当的 journal_mode。

3. **文件权限问题**
   
   问题：无法创建或访问数据库文件
   
   解决方案：检查文件系统权限，确保应用有适当的读写权限。