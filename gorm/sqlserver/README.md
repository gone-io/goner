# Gone Gorm SQL Server Driver

## Introduction

Gone Gorm SQL Server Driver is the Microsoft SQL Server database driver implementation for the Gone Gorm component. It allows you to use GORM to operate SQL Server databases in Gone applications, providing seamless integration with the Gone framework.

## Features

- Support for basic SQL Server database operations
- Seamless integration with Gone framework
- Connection pool configuration support
- Flexible logging configuration
- Transaction management support
- Database migration support

## Installation

```go
// Import Gone Gorm SQL Server driver in your application
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/sqlserver"
)

// Load SQL Server driver during application initialization
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // Load Gorm core component
            sqlserver.Load,   // Load SQL Server driver
            // ...
        ).
        Run()
}
```

## Configuration

### SQL Server Configuration

```properties
# SQL Server Basic Configuration
gorm.sqlserver.driver-name=              # Driver name, optional
gorm.sqlserver.dsn=sqlserver://gorm:gorm@localhost:9930?database=gorm  # Data source name
gorm.sqlserver.schema=                   # Database schema, default is dbo
gorm.sqlserver.encrypt=disable           # Connection encryption, options: disable, true, false
gorm.sqlserver.trust-server-certificate=false  # Whether to trust server certificate, default is false
gorm.sqlserver.app-name=                 # Application name, optional
```

## Usage Examples

### Basic Usage

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
    db *gorm.DB `gone:"*"` // Inject GORM instance
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

### Using Stored Procedures

```go
func (s *EmployeeService) GetEmployeesByDepartment(deptID int) ([]Employee, error) {
    var employees []Employee
    err := s.db.Raw("EXEC GetEmployeesByDepartment ?", deptID).Scan(&employees).Error
    return employees, err
}
```

## Best Practices

1. Database Design
   - Use appropriate data types
   - Design indexes properly
   - Follow SQL Server naming conventions
   - Use proper schema management

2. Connection Management
   - Configure connection pool parameters appropriately
   - Use connection encryption to protect data
   - Monitor connection status

3. Performance Optimization
   - Use appropriate indexes
   - Optimize query statements
   - Use stored procedures wisely
   - Avoid overuse of ORM features

4. Security
   - Apply principle of least privilege
   - Enable connection encryption
   - Update passwords regularly
   - Prevent SQL injection

## Common Issues

1. **Connection Issues**
   
   Issue: Unable to connect to SQL Server
   
   Solution: Check network connectivity, authentication information, and firewall settings, ensure DSN format is correct.

2. **Performance Issues**
   
   Issue: Slow query execution
   
   Solution: Check index usage, optimize query statements, use query plan analyzer.

3. **Character Set Issues**
   
   Issue: Special characters display as garbled text
   
   Solution: Ensure correct database and connection character set configuration, use NVARCHAR type for Unicode characters.