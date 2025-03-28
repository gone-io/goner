# Gone Gorm Component

## Introduction

Gone Gorm is the ORM component of the Gone framework, implemented based on [GORM](https://gorm.io/). It provides seamless integration with the Gone framework. Through this component, you can easily perform database operations in Gone applications using GORM, supporting various databases including MySQL, PostgreSQL, SQLite, SQL Server, and ClickHouse.

## Features

- Support for multiple databases: MySQL, PostgreSQL, SQLite, SQL Server, and ClickHouse
- Seamless integration with Gone framework
- Connection pool configuration support
- Flexible logging configuration
- Transaction management support
- Database migration support

## Database Driver Documentation

| Database Type | Documentation Link |
|--------------|-------------------|
| MySQL | [MySQL Driver Documentation](mysql/README.md) |
| PostgreSQL | [PostgreSQL Driver Documentation](postgres/README.md) |
| SQLite | [SQLite Driver Documentation](sqlite/README.md) |
| SQL Server | [SQL Server Driver Documentation](sqlserver/README.md) |
| ClickHouse | [ClickHouse Driver Documentation](clickhouse/README.md) |

## Installation

```go
// Import Gone Gorm component in your application
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    // Import specific database drivers as needed
    _ "github.com/gone-io/goner/gorm/mysql"
    // _ "github.com/gone-io/goner/gorm/postgres"
    // _ "github.com/gone-io/goner/gorm/sqlite"
    // _ "github.com/gone-io/goner/gorm/sqlserver"
    // _ "github.com/gone-io/goner/gorm/clickhouse"
)

// Load Gorm component during application initialization
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // Load Gorm core component
            mysql.Load, // Load MySQL driver
            // ...
        ).
        Run()
}
```

## Configuration

### Basic Configuration

```properties
# GORM Basic Configuration
gorm.skip-default-transaction=false     # Whether to skip default transaction, default is false
gorm.full-save-associations=false       # Whether to save associations completely, default is false
gorm.dry-run=false                      # Whether to generate SQL without executing, default is false
gorm.prepare-stmt=false                 # Whether to use prepared statements, default is false
gorm.disable-automatic-ping=false       # Whether to disable automatic ping, default is false
gorm.disable-foreign-key-constraint-when-migrating=false  # Whether to disable foreign key constraints during migration, default is false
gorm.ignore-relationships-when-migrating=false           # Whether to ignore relationships during migration, default is false
gorm.disable-nested-transaction=false    # Whether to disable nested transactions, default is false
gorm.allow-global-update=false           # Whether to allow global updates, default is false
gorm.query-fields=false                  # Whether to query all fields, default is false
gorm.create-batch-size=0                 # Batch creation size, default is 0
gorm.translate-error=false               # Whether to translate errors, default is false
gorm.propagate-unscoped=false            # Whether to propagate Unscoped to nested statements, default is false

# Connection Pool Configuration
gorm.pool.max-idle=10                    # Maximum number of idle connections, default is 10
gorm.pool.max-open=100                   # Maximum number of open connections, default is 100
gorm.pool.conn-max-lifetime=1h           # Maximum connection lifetime, default is 1 hour

# Logger Configuration
gorm.logger.slow-threshold=200ms         # Slow query threshold, default is 200ms
```

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

### PostgreSQL Configuration

```properties
# PostgreSQL Basic Configuration
gorm.postgres.driver-name=               # Driver name, optional
gorm.postgres.dsn=host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai  # Data source name
gorm.postgres.without-quoting-check=false  # Whether to skip quote checking, default is false
gorm.postgres.prefer-simple-protocol=false  # Whether to prefer simple protocol, default is false
gorm.postgres.without-returning=false    # Whether to disable RETURNING, default is false
```

### SQLite Configuration

```properties
# SQLite Basic Configuration
gorm.sqlite.driver-name=                 # Driver name, optional
gorm.sqlite.dsn=gorm.db                  # Data source name, default is gorm.db
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
    db *gorm.DB `gone:"*"`  // Inject GORM instance
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

### Auto Migration

```go
type AppStart struct {
    gone.Flag
    db *gorm.DB `gone:"*"`
}

func (s *AppStart) AfterRevive() error {
    // Auto migrate database structure
    return s.db.AutoMigrate(
        &User{},
        // Other models...
    )
}
```

### Transaction Handling

```go
func (s *UserService) TransferMoney(fromID, toID uint, amount float64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Execute database operations within transaction
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

1. Database Connection Management
   - Configure connection pool parameters appropriately to avoid connection leaks
   - Properly close database connections when the application shuts down
   - Use transactions for operations requiring atomicity

2. Model Design
   - Use struct tags to define field properties
   - Set up indexes and constraints appropriately
   - Use hook functions to handle model lifecycle events

3. Query Optimization
   - Use indexes to optimize query performance
   - Avoid N+1 query problems
   - Use preloading to reduce the number of queries

4. Logging Configuration
   - Set appropriate slow query thresholds
   - Enable detailed logging in development environment
   - Adjust log levels appropriately in production environment

## Common Issues

1. **Connection Pool Configuration**
   
   Issue: Application performance degradation or connection errors
   
   Solution: Adjust connection pool parameters like `gorm.pool.max-idle` and `gorm.pool.max-open` to match database server configuration.

2. **Slow Query Issues**
   
   Issue: Some queries take too long to execute
   
   Solution: Set `gorm.logger.slow-threshold` to identify slow queries, then optimize them by adding indexes or rewriting query logic.

3. **Transaction Handling**
   
   Issue: Improper error handling in transactions leading to data inconsistency
   
   Solution: Ensure proper error handling within transactions and roll back transactions when errors occur.