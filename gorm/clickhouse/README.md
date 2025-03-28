# Gone Gorm ClickHouse Driver

## Introduction

Gone Gorm ClickHouse Driver is the ClickHouse database driver implementation for the Gone Gorm component. It allows you to use GORM to operate ClickHouse databases in Gone applications, providing seamless integration with the Gone framework.

## Features

- Support for basic ClickHouse database operations
- Seamless integration with Gone framework
- Connection pool configuration support
- Flexible logging configuration

## Installation

```go
// Import Gone Gorm ClickHouse driver in your application
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/clickhouse"
)

// Load ClickHouse driver during application initialization
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // Load Gorm core component
            clickhouse.Load,  // Load ClickHouse driver
            // ...
        ).
        Run()
}
```

## Configuration

### ClickHouse Configuration

```properties
# ClickHouse Basic Configuration
gorm.clickhouse.driver-name=                  # Driver name, optional
gorm.clickhouse.dsn=tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20  # Data source name
gorm.clickhouse.debug=false                   # Whether to enable debug mode, default is false
gorm.clickhouse.cluster=                      # Cluster name, optional
gorm.clickhouse.settings=                     # Other ClickHouse specific settings, optional
```

## Usage Examples

### Basic Usage

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
    db *gorm.DB `gone:"*"` // Inject GORM instance
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

### Auto Migration

```go
type AppStart struct {
    gone.Flag
    db *gorm.DB `gone:"*"`
}

func (s *AppStart) AfterRevive() error {
    // Auto migrate database structure
    return s.db.AutoMigrate(
        &LogEntry{},
        // Other models...
    )
}
```

## Best Practices

1. Database Design
   - Use ClickHouse column types appropriately
   - Choose suitable table engines based on query patterns
   - Use appropriate partition keys and sorting keys

2. Query Optimization
   - Leverage ClickHouse's columnar storage features
   - Avoid JOIN operations
   - Use materialized views wisely

3. Write Optimization
   - Use batch inserts for better performance
   - Avoid frequent small batch writes
   - Set appropriate write timeout values

4. Monitoring and Maintenance
   - Monitor query performance and resource usage
   - Optimize table structure periodically
   - Clean up expired data timely

## Common Issues

1. **Connection Issues**
   
   Issue: Unable to connect to ClickHouse server
   
   Solution: Check network connectivity, port configuration, and authentication information, ensure DSN format is correct.

2. **Write Performance**
   
   Issue: Slow data write speed
   
   Solution: Use batch inserts, adjust write buffer size, choose appropriate table engine.

3. **Query Timeout**
   
   Issue: Complex queries timing out
   
   Solution: Optimize query statements, adjust timeout settings, use appropriate indexes and materialized views.