# goner/xorm

A powerful database integration component for the Gone framework, providing enhanced features based on XORM.

## Features

- Multiple database support (MySQL, SQLite3, PostgreSQL, Oracle, MSSQL)
- Automatic transaction management
- Transaction propagation
- Named parameter support
- Enhanced SQL query capabilities
- Master-slave cluster support
- Connection pool management

## Installation


```bash
go get github.com/gone-io/goner/xorm
```


## Import the required database driver:

```go
import (
    github.com/gone-io/goner/xorm

    // MySQL driver
    _ "github.com/go-sql-driver/mysql"

    // SQLite3 driver
    // _ "github.com/mattn/go-sqlite3"

    // PostgreSQL driver
    // _ "github.com/lib/pq"

    // Oracle driver
    // _ "github.com/mattn/go-oci8"

    // MSSQL driver
    // _ "github.com/denisenkom/go-mssqldb"
)
```




## Configuration

### Single Database Mode

```ini
database.driver-name=mysql
database.dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

### Cluster Mode

```ini
database.cluster.enable=true

# Master database configuration
database.cluster.master.driver-name=mysql
database.cluster.master.dsn=root:123456@tcp(master-db-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local

# Slave database configuration
database.cluster.slaves[0].driver-name=mysql
database.cluster.slaves[0].dsn=root:123456@tcp(slave-db-0-host:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

### Configuration Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| database.cluster.enable | No | false | Enable cluster mode |
| database.driver-name | No* | - | Database driver name |
| database.dsn | No* | - | Database connection string |
| database.max-idle-count | No | 5 | Maximum idle connections |
| database.max-open | No | 20 | Maximum open connections |
| database.max-lifetime | No | 10m | Connection maximum lifetime |
| database.show-sql | No | true | Show SQL logs |

*Required in non-cluster mode

##  Load And Start
```go
package main
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/gone/goner/xorm"
    _ "github.com/go-sql-driver/mysql"
)
func main() {
    gone.
        NewApp(
            xorm.Load,
        ).
        Run(func(db xorm.Engine){
            // Your application logic
            db.Get(/*...*/)
        })
}
```
## Usage

### Basic Usage

```go
type dbUser struct {
    gone.Flag
    db xorm.Engine `gone:"*"` // Inject database engine
}

type Book struct {
    Id    int64
    Title string
}

func (d *dbUser) GetBookById(id int64) (book *Book, err error) {
    book = new(Book)
    has, err := d.db.Where("id=?", id).Get(book)
    if err != nil {
        return nil, gone.ToError(err)
    }
    if !has {
        return nil, gone.NewParameterError("book not found", 404)
    }
    return book, nil
}
```

### Automatic Transaction

```go
func (d *db) updateUser(user *entity.User) error {
    return d.Transaction(func(session xorm.Interface) error {
        _, err := session.ID(user.Id).Update(user)
        return gone.ToError(err)
    })
}
```

### Named Parameters

```go
sql, args := xorm.MustNamed(`
    update user
    set status = :status,
        avatar = :avatar
    where id = :id`,
    map[string]any{
        "id":     1,
        "status": 1,
        "avatar": "https://example.com/avatar.jpg",
    },
)
```

## Best Practices

1. Connection Pool Configuration
   - Adjust `max-idle-count` based on load
   - Set appropriate `max-open` to prevent database overload
   - Configure `max-lifetime` to manage connection lifecycle

2. Cluster Mode Optimization
   - Use read/write separation effectively
   - Configure connection pools for master and slave databases separately
   - Consider master-slave replication lag in transactions

3. SQL Optimization
   - Use `show-sql` in development for SQL debugging
   - Implement proper indexing
   - Avoid large transactions
   - Use batch operations instead of loops

## License

MIT License