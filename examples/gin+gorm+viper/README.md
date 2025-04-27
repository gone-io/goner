[//]: # (desc: simple web demo using gin, gorm, viper, mysql)
<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Framework Integration Example with Gin+GORM+Viper

This example demonstrates how to integrate the Gone framework with Gin, GORM, and Viper components to create a simple web application.

## Project Overview

This example demonstrates the following features:

- Using Gone framework's dependency injection mechanism
- Integrating Gin framework for web routing and controllers
- Integrating GORM framework for database access
- Using Viper for configuration management

## Project Structure

```
.
├── config/
│   └── default.properties  # Configuration file
├── go.mod                  # Go module definition
└── main.go                 # Main program
```

## Configuration

The configuration file `config/default.properties` contains MySQL database connection information:

```properties
gorm.mysql.dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

## Code Implementation

### Main Program

`main.go` contains the core logic of the application:

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
	goneGorm "github.com/gone-io/goner/gorm"
	"github.com/gone-io/goner/gorm/mysql"
	"gorm.io/gorm"
)

// Define controller
type HelloController struct {
	gone.Flag
	gin.IRouter `gone:"*"`      // Inject router
	uR          *UserRepository `gone:"*"`
}

// Mount implements gin.Controller interface
func (h *HelloController) Mount() gin.MountError {
	h.GET("/hello", h.hello) // Register route
	h.GET("/user/:id", h.getUser)
	return nil
}

func (h *HelloController) hello() (string, error) {
	return "Hello, Gone!", nil
}
func (h *HelloController) getUser(in struct {
	id uint `param:"id"`
}) (*User, error) {

	user, err := h.uR.GetByID(in.id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Define data model and repository
type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type UserRepository struct {
	gone.Flag
	*gorm.DB `gone:"*"`
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	var user User
	err := r.First(&user, id).Error
	return &user, err
}

func main() {
	// Load components and start application
	gone.
		Loads(
			goner.BaseLoad,
			goneGorm.Load, // Load Gorm core component
			mysql.Load,    // Load MySQL driver
			gin.Load,      // Load Gin component
		).
		Load(&HelloController{}). // Load controller
		Load(&UserRepository{}).  // Load repository
		Serve()
}
```

## Code Analysis

### Dependency Injection

Gone framework uses the `gone:"*"` tag for dependency injection:

1. `gin.IRouter` and `UserRepository` are injected into `HelloController`
2. `*gorm.DB` is injected into `UserRepository`

### Controller

`HelloController` implements the `gin.Controller` interface's `Mount` method and registers two routes:

- `GET /hello`: Returns a simple greeting message
- `GET /user/:id`: Queries user information by ID

### Data Model and Repository

- `User` struct defines the user model
- `UserRepository` provides database access methods, such as `GetByID`

### Application Startup

In the `main` function:

1. Use `gone.Loads` to load basic components:
   - `goner.BaseLoad`: Basic component
   - `goneGorm.Load`: GORM core component
   - `mysql.Load`: MySQL driver
   - `gin.Load`: Gin component

2. Load custom components:
   - `&HelloController{}`: Controller
   - `&UserRepository{}`: Repository

3. Call `Serve()` to start the application

## Running the Example

### Environment Setup

1. Ensure Go environment is installed (Go 1.16+ recommended)
2. Prepare MySQL database, create `test` database
3. Modify database connection information in `config/default.properties` as needed

### Create Database Table

Execute the following SQL in MySQL to create the users table:

```sql
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Insert test data
INSERT INTO `users` (`id`, `name`) VALUES (1, 'Gone User');
```

### Start Application

```bash
go run main.go
```

### Test API

1. Access greeting endpoint:
   ```
   curl http://localhost:8080/hello
   ```
   Expected response: `"Hello, Gone!"`

2. Query user information:
   ```
   curl http://localhost:8080/user/1
   ```
   Expected response: `{"ID":1,"Name":"Gone User"}`

## Summary

This example demonstrates how to integrate the Gone framework with Gin, GORM, and Viper to create a simple but fully functional web application. Through Gone framework's dependency injection mechanism, the coupling between components is reduced, making the code clearer and easier to maintain.

This example can serve as a starting point for developing web applications using the Gone framework. You can extend it with more features, such as adding middleware, implementing more complex business logic, and more.