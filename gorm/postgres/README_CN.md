<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/gorm/postgres 组件， Gone Gorm PostgreSQL 驱动

## 简介

Gone Gorm PostgreSQL 驱动是 Gone Gorm 组件的 PostgreSQL 数据库驱动实现。它允许您在 Gone 应用中使用 GORM 操作 PostgreSQL 数据库，提供了与 Gone 框架的无缝集成。

## 特性

- 支持 PostgreSQL 数据库的基本操作
- 与 Gone 框架无缝集成
- 支持连接池配置
- 提供灵活的日志配置
- 支持事务管理
- 支持数据库迁移
- 支持 JSON 和 JSONB 数据类型
- 支持数组类型

## 安装

```go
// 在您的应用中引入 Gone Gorm PostgreSQL 驱动
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/postgres"
)

// 在应用初始化时加载 PostgreSQL 驱动
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // 加载 Gorm 核心组件
            postgres.Load,    // 加载 PostgreSQL 驱动
            // ...
        ).
        Run()
}
```

## 配置说明

### PostgreSQL 配置

```properties
# PostgreSQL 基础配置
gorm.postgres.driver-name=               # 驱动名称，可选
gorm.postgres.dsn=host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai  # 数据源名称
gorm.postgres.without-quoting-check=false  # 是否不进行引号检查，默认为 false
gorm.postgres.prefer-simple-protocol=false  # 是否优先使用简单协议，默认为 false
gorm.postgres.without-returning=false    # 是否不使用 RETURNING，默认为 false
```

## 使用示例

### 基本用法

```go
package example

import (
    "github.com/gone-io/gone/v2"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint            `gorm:"primaryKey"`
    Name        string          `gorm:"size:255"`
    Description string          `gorm:"type:text"`
    Price       float64         `gorm:"type:decimal(10,2)"`
    Tags        []string        `gorm:"type:text[]"`
    Metadata    map[string]any  `gorm:"type:jsonb"`
}

type ProductService struct {
    gone.Flag
    db *gorm.DB `gone:"*"` // 注入 GORM 实例
}

func (s *ProductService) CreateProduct(product *Product) error {
    return s.db.Create(product).Error
}

func (s *ProductService) GetProductByID(id uint) (*Product, error) {
    var product Product
    if err := s.db.First(&product, id).Error; err != nil {
        return nil, err
    }
    return &product, nil
}

func (s *ProductService) SearchProducts(query string) ([]Product, error) {
    var products []Product
    if err := s.db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&products).Error; err != nil {
        return nil, err
    }
    return products, nil
}
```

### 使用 JSON 和数组类型

```go
func (s *ProductService) UpdateProductMetadata(id uint, metadata map[string]any) error {
    return s.db.Model(&Product{}).Where("id = ?", id).Update("metadata", metadata).Error
}

func (s *ProductService) AddProductTag(id uint, tag string) error {
    return s.db.Model(&Product{}).Where("id = ?", id).Update("tags", gorm.Expr("array_append(tags, ?)", tag)).Error
}

func (s *ProductService) FindProductsByTags(tags []string) ([]Product, error) {
    var products []Product
    if err := s.db.Where("tags && ?", tags).Find(&products).Error; err != nil {
        return nil, err
    }
    return products, nil
}
```

## 最佳实践

1. 数据库设计
   - 合理使用 PostgreSQL 特有的数据类型
   - 适当使用索引提升查询性能
   - 利用 JSONB 类型存储非结构化数据
   - 合理使用数组类型

2. 连接管理
   - 配置适当的连接池参数
   - 设置合理的超时时间
   - 使用 SSL 加密保护数据传输

3. 查询优化
   - 使用适当的索引
   - 优化复杂查询
   - 利用 PostgreSQL 的全文搜索功能
   - 合理使用事务

4. 性能优化
   - 使用批量操作
   - 合理设置 vacuum 策略
   - 监控查询性能
   - 定期维护数据库

## 常见问题

1. **连接问题**
   
   问题：无法连接到 PostgreSQL 服务器
   
   解决方案：检查网络连接、端口配置和认证信息，确保 DSN 格式正确。

2. **性能问题**
   
   问题：查询执行缓慢
   
   解决方案：检查索引使用情况，优化查询语句，使用 EXPLAIN ANALYZE 分析查询计划。

3. **JSON 查询问题**
   
   问题：JSONB 字段查询效率低
   
   解决方案：为 JSONB 字段创建适当的索引，优化查询语句。