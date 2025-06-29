<p align="center">
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/mongo component and Gone MongoDB Integration

This package provides MongoDB integration for Gone applications, offering easy-to-use client configuration and management.

## Features

- Easy integration with Gone's dependency injection system
- Support for multiple MongoDB client instances
- Singleton client instance management
- Comprehensive configuration options
- Connection pooling and timeout management
- Authentication support

## Installation

```bash
gonectl install goner/mongo
```

## Configuration

Create a `default.yaml` file in your project's config directory with the following MongoDB configuration:

```yaml
mongo:
  uri: "mongodb://localhost:27017"     # MongoDB connection URI
  database: "myapp"                    # Default database name
  username: ""                         # Optional: Username for authentication
  password: ""                         # Optional: Password for authentication
  authSource: "admin"                  # Optional: Authentication database
  maxPoolSize: 100                     # Optional: Maximum connection pool size
  minPoolSize: 0                       # Optional: Minimum connection pool size
  maxConnIdleTime: "30m"               # Optional: Maximum connection idle time
  connectTimeout: "10s"                # Optional: Connection timeout
  socketTimeout: "30s"                 # Optional: Socket timeout
  serverSelectionTimeout: "30s"        # Optional: Server selection timeout
```

### Configuration Options

- **uri**: MongoDB connection string. Can include host, port, database, and other options
- **database**: Default database name to use
- **username/password**: Credentials for authentication
- **authSource**: Database to authenticate against (default: "admin")
- **maxPoolSize**: Maximum number of connections in the pool
- **minPoolSize**: Minimum number of connections in the pool
- **maxConnIdleTime**: Maximum time a connection can remain idle
- **connectTimeout**: Timeout for establishing connections
- **socketTimeout**: Timeout for socket operations
- **serverSelectionTimeout**: Timeout for server selection

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/mongo"
    "go.mongodb.org/mongo-driver/bson"
    mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
    gone.Flag
    mongoClient *mongoDriver.Client `gone:"*"`
    logger      gone.Logger          `gone:"*"`
}

func (s *UserService) CreateUser(name, email string) error {
    collection := s.mongoClient.Database("myapp").Collection("users")
    
    user := bson.M{
        "name":  name,
        "email": email,
    }
    
    _, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        s.logger.Errorf("Failed to create user: %v", err)
        return err
    }
    
    s.logger.Infof("User created successfully: %s", name)
    return nil
}

func (s *UserService) GetUser(email string) (bson.M, error) {
    collection := s.mongoClient.Database("myapp").Collection("users")
    
    var user bson.M
    err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        s.logger.Errorf("Failed to get user: %v", err)
        return nil, err
    }
    
    return user, nil
}

func main() {
    gone.NewApp(
        gone.Load(mongo.Load),
        gone.Load(&UserService{}),
    ).Run()
}
```

### Multiple Database Connections

You can configure multiple MongoDB connections:

```yaml
mongo:
  uri: "mongodb://localhost:27017"
  database: "main"

mongo-analytics:
  uri: "mongodb://analytics-server:27017"
  database: "analytics"
  username: "analytics_user"
  password: "analytics_pass"
```

```go
type AnalyticsService struct {
    gone.Flag
    mainClient      *mongoDriver.Client `gone:"*"`
    analyticsClient *mongoDriver.Client `gone:"*,mongo-analytics"`
}
```

## Error Handling

The component includes comprehensive error handling:

- Connection failures are properly reported
- Configuration errors are detailed
- Connection pooling errors are handled gracefully

## Best Practices

1. **Connection Pooling**: Configure appropriate pool sizes based on your application's needs
2. **Timeouts**: Set reasonable timeouts to avoid hanging operations
3. **Authentication**: Use authentication in production environments
4. **Database Selection**: Specify the database name in configuration for clarity
5. **Error Handling**: Always handle errors returned by MongoDB operations

## Dependencies

- [go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver) - Official MongoDB Go driver
- [github.com/gone-io/gone/v2](https://github.com/gone-io/gone) - Gone framework