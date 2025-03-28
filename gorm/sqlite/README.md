# Gone Gorm SQLite Driver

## Introduction

Gone Gorm SQLite Driver is the SQLite database driver implementation for the Gone Gorm component. It allows you to use GORM to operate SQLite databases in Gone applications, providing seamless integration with the Gone framework. SQLite is a lightweight, zero-configuration, self-contained database engine that is particularly suitable for embedded applications and development environments.

## Features

- Support for basic SQLite database operations
- Seamless integration with Gone framework
- Zero configuration, easy to use
- Transaction management support
- Database migration support
- Suitable for development and testing environments
- In-memory database support

## Installation

```go
// Import Gone Gorm SQLite driver in your application
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/gorm"
    _ "github.com/gone-io/goner/gorm/sqlite"
)

// Load SQLite driver during application initialization
func main() {
    gone.
        Loads(
            // ...
            gorm.Load,        // Load Gorm core component
            sqlite.Load,      // Load SQLite driver
            // ...
        ).
        Run()
}
```

## Configuration

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

type Note struct {
    ID      uint   `gorm:"primaryKey"`
    Title   string `gorm:"size:255"`
    Content string `gorm:"type:text"`
}

type NoteService struct {
    gone.Flag
    db *gorm.DB `gone:"*"` // Inject GORM instance
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

### Using In-Memory Database

```go
// Configure in-memory database
gorm.sqlite.dsn=:memory:

// Or use directly in code
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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
        &Note{},
        // Other models...
    )
}
```

## Best Practices

1. Database Design
   - Design table structure appropriately
   - Use indexes properly
   - Avoid overly complex relationships

2. Performance Optimization
   - Use transactions for batch operations
   - Perform VACUUM operations periodically
   - Configure journal_mode appropriately

3. Concurrency Handling
   - Be aware of SQLite's concurrency limitations
   - Use appropriate transaction isolation levels
   - Avoid long-lasting locks

4. Backup and Maintenance
   - Backup database files regularly
   - Monitor database size
   - Clean up unnecessary data timely

## Common Issues

1. **Concurrency Access Issues**
   
   Issue: Database locks due to multiple concurrent connections
   
   Solution: Use appropriate locking strategies, avoid long transactions, consider using WAL mode.

2. **Performance Issues**
   
   Issue: Database operations becoming slow
   
   Solution: Run VACUUM periodically, optimize indexes, use appropriate journal_mode.

3. **File Permission Issues**
   
   Issue: Unable to create or access database file
   
   Solution: Check filesystem permissions, ensure application has appropriate read/write permissions.