<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/viper Component

The **goner/viper** component is the configuration management component of the Gone framework, implemented based on [spf13/viper](https://github.com/spf13/viper). It provides flexible and powerful configuration management capabilities for your Gone applications, supporting multiple configuration sources and formats, making your application configuration simple and efficient.

## Features

- Seamless integration with Gone framework
- Multiple configuration sources: files, environment variables, command-line arguments
- Multiple configuration file formats: JSON, YAML, TOML, Properties
- Hierarchical configuration structure and default value mechanism
- Environment variable override capability for enhanced flexibility

## Configuration File Search Mechanism

The component automatically searches for configuration files in the following priority order:

1. Executable file directory
2. `config` subdirectory under the executable file directory
3. Current working directory
4. `config` subdirectory under the current working directory
5. Additional paths when starting Gone with the `Test` function:
    - `config` directory under the go.mod file directory
    - `testdata` directory under the current working directory
    - `testdata/config` directory under the current working directory
6. If the `CONF` environment variable is set or the `-conf` option is used at startup, it will also search the specified configuration file path

## Configuration File Loading Order

When multiple configuration files exist in the same directory, the component loads and merges configurations in the following order:

### Default Configuration Files (in order)
1. "default.json"
2. "default.toml"
3. "default.yaml"
4. "default.yml"
5. "default.properties"

### Environment-Specific Configuration Files
The default environment is `local`, which can be modified through the `ENV` environment variable or `-env` startup parameter. Environment configuration files are loaded in the following order:

1. "${env}.json"
2. "${env}.toml"
3. "${env}.yaml"
4. "${env}.yml"
5. "${env}.properties"

### Test-Specific Configuration Files
When starting Gone using the `Test` function, additional test-specific configuration files `${default|env}_test.${ext}` will be loaded.

For example, in the following test code:

```go
func TestCase(t *testing.T){
    gone.
        Loads(
            viper.Load, // Load configuration component
            // Other components...
        ).
        Test(func(){
            // Test code
        })
}
```

The system will load existing configuration files in sequence (non-existent files will be ignored):
1. "default.json"
2. "default.toml"
3. "default.yaml"
4. "default.yml"
5. "default.properties"
6. "default_test.json"
7. "default_test.toml"
8. "default_test.yaml"
9. "default_test.yml"
10. "default_test.properties"
11. "local.json"
12. "local.toml"
13. "local.yaml"
14. "local.yml"
15. "local.properties"
16. "local_test.json"
17. "local_test.toml"
18. "local_test.yaml"
19. "local_test.yml"
20. "local_test.properties"

**Important Notes:**
1. Content from multiple configuration files will be automatically merged, with later loaded configurations overriding earlier ones for the same keys
2. Variable substitution in properties files is only valid within the same configuration file, cross-file substitution is not supported

## Quick Start

### 1. Install the Component

```bash
go install github.com/gone-io/goner/viper
```

### 2. Load Configuration Component in Your Application

```go
package main

import (
    "github.com/gone-io/v2"
    "github.com/gone-io/goner/viper"
)

func main() {
    gone.
        Loads(
            viper.Load, // Load configuration component
            // Other components...
        ).
        Run() // Or use Serve()
}
```

### 3. Inject Configuration via Tags

```go
type MyService struct {
    gone.Flag
    
    // Inject configuration via gone:"config" tag, supports default values
    ServerHost string `gone:"config,server.host,default=localhost"`
    ServerPort int    `gone:"config,server.port,default=8080"`
    DbURL      string `gone:"config,db.url"`
}

func (s *MyService) Start() error {
    // Use injected configuration values
    fmt.Printf("Service running at %s:%d\n", s.ServerHost, s.ServerPort)
    return nil
}
```

### 4. Manually Get Configuration Values

```go
type MyComponent struct {
    gone.Flag
    conf gone.Configure `gone:"*"` // Inject configuration manager
}

func (c *MyComponent) DoSomething() error {
    // Get string configuration
    var host string
    err := c.conf.Get("server.host", &host, "localhost")
    if err != nil {
        return err
    }
    
    // Get integer configuration
    var port int
    err = c.conf.Get("server.port", &port, "8080")
    if err != nil {
        return err
    }
    
    // Get complex struct configuration
    var dbConfig struct {
        URL      string
        Username string
        Password string
    }
    err = c.conf.Get("db", &dbConfig, "")
    if err != nil {
        return err
    }
    
    return nil
}
```

## Configuration File Format Examples

### Properties Format

```properties
# Server configuration
server.host=localhost
server.port=8080

# Database configuration
db.username=root
db.password=secret
db.database=mydb

# Variable substitution supported within the same file
db.url=mysql://localhost:3306/${db.database}

# Log configuration
log.level=info
log.path=/var/log/myapp
```

### YAML Format

```yaml
server:
  host: localhost
  port: 8080

db:
  url: mysql://localhost:3306/mydb
  username: root
  password: secret

log:
  level: info
  path: /var/log/myapp
```

## Using Environment Variables to Override Configuration

You can override values in configuration files using environment variables for increased deployment flexibility. Environment variable naming follows the rule `GONE_config_key`, where dots (`.`) in the config key are replaced with underscores (`_`).

For example, to override the `server.port` configuration, set the environment variable:
```
GONE_SERVER_PORT=9090
```

## API Reference

### Configure Interface

```go
type Configure interface {
    Get(key string, v any, defaultVal string) error
}
```

Core interface for getting configuration values, parameter description:

- `key`: Configuration key name, supports dot-separated hierarchical structure
- `v`: Pointer to variable for receiving configuration value
- `defaultVal`: Default value used when configuration doesn't exist

## Best Practice Recommendations

1. Use hierarchical structure to organize configurations for better readability and maintainability
2. Provide reasonable default values for critical configurations to ensure application runs normally when configurations are missing
3. Inject sensitive information (like passwords, API keys) through environment variables, avoid hardcoding in configuration files
4. Use different environment configuration files (like `dev.yaml`, `prod.yaml`) to manage configurations for different environments
5. Use lowercase letters and dots for configuration key names, maintain consistent naming conventions