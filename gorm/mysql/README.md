# Gone Gorm MySQL Driver

## Introduction

Gone Gorm MySQL Driver is the MySQL database driver implementation for the Gone Gorm component. It allows you to use GORM to operate MySQL databases in Gone applications, providing seamless integration with the Gone framework.

## Features

- Support for basic MySQL database operations
- Seamless integration with Gone framework
- Connection pool configuration support
- Flexible logging configuration
- Transaction management support
- Database migration support

## Installation

```go
// Import Gone Gorm MySQL driver in your application
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/mysql"
)

// Load MySQL driver during application initialization
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // Load Gorm core component
            mysql.Load,       // Load MySQL driver
            // ...
        ).
        Run()
}
```

## Configuration

### MySQL Configuration

```properties
# MySQL Basic Configuration
gorm.mysql.driver-name=                  # Driver name, optional
gorm.mysql.dsn=user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local  # Data source name
gorm.mysql.server-version=               # Server version, optional
gorm.mysql.skip-initialize-with-version=false  # Whether to skip initialization with version, default is false
gorm.mysql.default-string-size=0         # Default string size, default is 0
gorm.mysql.default-datetime-precision=   # Default datetime precision, optional
gorm.mysql.disable-with-returning=false  # Whether to disable WITH RETURNING, default is false
gorm.mysql.disable-datetime-precision=false  # Whether to disable datetime precision, default is false
gorm.mysql.dont-support-rename-index=false  # Whether rename index is not supported, default is false
gorm.mysql.dont-support-rename-column=false  # Whether rename column is not supported, default is false
gorm.mysql.dont-support-for-share-clause=false  # Whether FOR SHARE clause is not supported, default is false
gorm.mysql.dont-support-null-as-default-value=false  # Whether NULL as default value is not supported, default is false
gorm.mysql.dont-support-rename-column-unique=false  # Whether rename unique column is not supported, default is false
gorm.mysql.dont-support-drop-constraint=false  # Whether drop constraint is not supported, default is false
```

## Usage Examples

### Basic Usage

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
    db *gorm.DB `gone:"*"` // Inject GORM instance
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

### Transaction Handling

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
            return errors.New("insufficient balance")
        }
        
        // Update account balances
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

## Best Practices

1. Database Design
   - Design table structure and field types appropriately
   - Use proper indexes to improve query performance
   - Follow database design principles

2. Connection Management
   - Configure connection pool parameters appropriately
   - Close unnecessary connections timely
   - Use connection pool monitoring tools

3. Query Optimization
   - Use appropriate indexes
   - Avoid full table scans
   - Optimize JOIN queries
   - Use prepared statements

4. Transaction Handling
   - Use appropriate transaction isolation levels
   - Avoid long transactions
   - Handle transaction rollbacks properly

## Common Issues

1. **Connection Issues**
   
   Issue: Unable to connect to MySQL server
   
   Solution: Check network connectivity, port configuration, and authentication information, ensure DSN format is correct.

2. **Character Set Issues**
   
   Issue: Chinese or special characters display as garbled text
   
   Solution: Ensure using utf8mb4 character set and configure charset parameter correctly in DSN.

3. **Performance Issues**
   
   Issue: Slow query execution
   
   Solution: Check index usage, optimize query statements, adjust database configuration parameters.