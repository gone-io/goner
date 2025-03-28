# Gone Gorm PostgreSQL Driver

## Introduction

Gone Gorm PostgreSQL Driver is the PostgreSQL database driver implementation for the Gone Gorm component. It allows you to use GORM to operate PostgreSQL databases in Gone applications, providing seamless integration with the Gone framework.

## Features

- Support for basic PostgreSQL database operations
- Seamless integration with Gone framework
- Connection pool configuration support
- Flexible logging configuration
- Transaction management support
- Database migration support
- JSON and JSONB data type support
- Array type support

## Installation

```go
// Import Gone Gorm PostgreSQL driver in your application
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/postgres"
)

// Load PostgreSQL driver during application initialization
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // Load Gorm core component
            postgres.Load,    // Load PostgreSQL driver
            // ...
        ).
        Run()
}
```

## Configuration

### PostgreSQL Configuration

```properties
# PostgreSQL Basic Configuration
gorm.postgres.driver-name=               # Driver name, optional
gorm.postgres.dsn=host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai  # Data source name
gorm.postgres.without-quoting-check=false  # Whether to skip quote checking, default is false
gorm.postgres.prefer-simple-protocol=false  # Whether to prefer simple protocol, default is false
gorm.postgres.without-returning=false    # Whether to disable RETURNING, default is false
```

## Usage Examples

### Basic Usage

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
    db *gorm.DB `gone:"*"` // Inject GORM instance
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

### Using JSON and Array Types

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

## Best Practices

1. Database Design
   - Use PostgreSQL-specific data types appropriately
   - Set up proper indexes for better query performance
   - Leverage JSONB type for unstructured data
   - Use array types effectively

2. Connection Management
   - Configure connection pool parameters appropriately
   - Set reasonable timeout values
   - Use SSL encryption for data transmission

3. Query Optimization
   - Use appropriate indexes
   - Optimize complex queries
   - Utilize PostgreSQL's full-text search capabilities
   - Use transactions wisely

4. Performance Optimization
   - Use batch operations
   - Set up proper vacuum strategies
   - Monitor query performance
   - Maintain database regularly

## Common Issues

1. **Connection Issues**
   
   Issue: Unable to connect to PostgreSQL server
   
   Solution: Check network connectivity, port configuration, and authentication information, ensure DSN format is correct.

2. **Performance Issues**
   
   Issue: Slow query execution
   
   Solution: Check index usage, optimize query statements, use EXPLAIN ANALYZE to analyze query plans.

3. **JSON Query Issues**
   
   Issue: Low efficiency in JSONB field queries
   
   Solution: Create appropriate indexes for JSONB fields, optimize query statements.