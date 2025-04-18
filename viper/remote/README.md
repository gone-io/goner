# Gone Viper Remote

## What is Gone Viper Remote?

The `remote` package is a crucial component of the Gone framework that extends the Viper configuration system, enabling your applications to fetch configuration information from remote configuration centers (such as etcd, consul, etc.). Built upon [spf13/viper/remote](https://github.com/spf13/viper/tree/master/remote), this package is specifically optimized for the Gone framework, providing seamless integration.

Imagine having multiple application instances that need to share the same configuration, or needing to dynamically update configurations without restarting your application — this is where Gone Viper Remote shines.

## Why Choose Gone Viper Remote?

- **Centralized Configuration Management**: All application instances can fetch the latest configuration from a single configuration center
- **Real-time Configuration Updates**: Supports hot configuration updates without application restart
- **Enhanced Security**: Supports encrypted configuration to protect sensitive information
- **Local Configuration Fallback**: Automatically falls back to local configuration when remote configuration is unavailable
- **Multiple Data Source Support**: Compatible with various popular configuration centers

## Getting Started

### Step 1: Install the Package

```bash
go get github.com/gone-io/goner/viper/remote
```

### Step 2: Load the Component in Your Application

```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/viper/remote"
)

func main() {
    // Create Gone application and load remote component
    gone.NewApp(remote.Load).Run()
}
```

### Step 3: Configure Remote Provider

Set up the remote configuration provider in your configuration file (e.g., `config/default.yaml`):

```yaml
# config/default.yaml
viper:
  remote:
    type: yaml
    watch: true                     # Enable hot configuration updates
    watchDuration: 5s               # Check for changes every 5 seconds
    useLocalConfIfKeyNotExist: true # Use local configuration if remote key not found
    providers:
      - provider: etcd              # Provider type
        endpoint: localhost:2379    # Provider address
        path: /config/myapp         # Configuration path
        configType: json            # Configuration format
        keyring:                    # Encryption key for configuration (optional)
      - provider: consul
        endpoint: localhost:8500
        path: myapp/config
        configType: yaml
        keyring:
```

## Enhanced Security: Using Encrypted Configuration

For sensitive information (such as database passwords, API keys), you can use encrypted configuration to enhance security.

### Setting up GPG Keys

1. Generate GPG key pair:

```bash
# Generate GPG key pair
gpg --gen-key

# Export public key (for encryption)
gpg --export > pubring.gpg

# Export private key (for decryption)
gpg --export-secret-keys > secring.gpg
```

2. Specify the key file in configuration:

```yaml
viper.remote:
  providers:
    - provider: etcd3
      endpoint: http://localhost:2379
      path: /config/secure-config
      configType: yaml
      keyring: /path/to/secring.gpg  # Specify key file path
```

### Example Using Encrypted Configuration

```go
package main

import (
    "fmt"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/viper/remote"
)

func main() {
    gone.
        NewApp(remote.Load).
        Run(func(params struct {
            apiKey string `gone:"config,secure.api.key"`
            dbPass string `gone:"config,secure.database.password"`
        }) {
            fmt.Printf("API Key: %s, DB Password: %s\n", params.apiKey, params.dbPass)
        })
}
```

## Configuration Details

### Provider Specification

Each remote provider is defined by the following attributes:

```go
type Provider struct {
    Provider   string // Provider type: etcd, consul, etc.
    Endpoint   string // Provider address
    Path       string // Configuration path in provider
    ConfigType string // Configuration format: json, yaml, etc.
    Keyring    string // Encryption key for configuration (optional)
}
```

### Global Configuration Options

| Option | Description | Default | Usage Recommendation |
|--------|-------------|---------|---------------------|
| viper.remote.providers | List of remote configuration providers | [] | Configure multiple providers for redundancy |
| viper.remote.watch | Enable hot configuration updates | false | Recommended for production environments |
| viper.remote.watchDuration | Interval for checking configuration updates | 5s | Adjust based on configuration change frequency |
| viper.remote.useLocalConfIfKeyNotExist | Use local configuration if remote key not found | true | Enable to improve system reliability |

## Supported Remote Providers

Currently supports the following remote configuration centers:

- **etcd/etcd3**: High-availability distributed key-value store, suitable for large-scale clusters
- **consul**: Service discovery and configuration tool with built-in health checks
- **firestore**: Google Cloud's NoSQL database
- **nats**: High-performance distributed messaging system

## Working Principle Analysis

Gone Viper Remote works as follows:

1. **Initialization**: Read remote provider information from local configuration file
2. **Connection**: Connect to remote configuration center
3. **Loading**: Fetch configuration information from remote and merge with local configuration
4. **Monitoring** (if enabled): Periodically check for remote configuration changes and update timely

### Local Configuration Fallback Mechanism

When the remote configuration center is unavailable or a configuration key doesn't exist, the system automatically falls back to local configuration:

1. Application requests configuration value
2. System first tries to fetch from remote
3. If fetch fails and `useLocalConfIfKeyNotExist` is `true`
4. System falls back to local configuration file
5. If local configuration doesn't exist, use default value (if provided)

This mechanism is particularly suitable for:

- **Development Environment**: Developers can override certain configurations locally
- **Disaster Recovery**: Application can still run when remote configuration center is unavailable
- **Configuration Migration**: Gradually migrate from local to remote configuration center

## Practical Example: Complete Application

Here's a complete example using etcd as configuration center:

```go
package main

import (
    "fmt"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/viper/remote"
    "time"
)

// Define database configuration structure
type Database struct {
    UserName string `mapstructure:"username"`
    Pass     string `mapstructure:"password"`
}

func main() {
    gone.
        NewApp(remote.Load).
        Run(func(params struct {
            serverName string   `gone:"config,server.name"`
            serverPort int      `gone:"config,server.port"`
            dbUserName string   `gone:"config,database.username"`
            dbUserPass string   `gone:"config,database.password"`
            database   *Database `gone:"config,database"`
            key        string   `gone:"config,key.not-existed-in-etcd"`
        }) {
            // Print configuration information
            fmt.Printf("Server Name: %s, Port: %d\n", params.serverName, params.serverPort)
            fmt.Printf("Database User: %s, Password: %s\n", params.dbUserName, params.dbUserPass)
            fmt.Printf("Local Config Item: %s\n", params.key)

            // Print database configuration every 10 seconds to demonstrate hot updates
            for i := 0; i < 10; i++ {
                fmt.Printf("Database Configuration: %#+v\n", *params.database)
                time.Sleep(10 * time.Second)
            }
        })
}
```

Configuration file setup:

```yaml
# config/default.yaml
viper.remote:
  type: yaml
  watch: true
  watchDuration: 5s
  useLocalConfIfKeyNotExist: true
  providers:
    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path: /config/application.yaml
    
    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path: /config/database.yaml

# Local configuration, used when remote key doesn't exist
key:
  not-existed-in-etcd: 1000
```

Configuration content in etcd:

```yaml
# /config/application.yaml
server.name: config-demo
server.port: 9090

# /config/database.yaml
database:
  username: config-demo
  password: config-demo-password
```

## Best Practice Recommendations

1. **Layered Configuration**: Separate configurations by functionality, store in different paths
2. **Regular Backups**: Regularly backup remote configuration center data
3. **Appropriate Monitoring Interval**: Adjust `watchDuration` based on application needs, avoid too frequent checks
4. **Encrypt Sensitive Information**: Use encryption for passwords, API keys, and other sensitive information
5. **Local Configuration Fallback**: Keep local configuration consistent with remote configuration as an emergency measure

## Frequently Asked Questions

1. **What to do when remote configuration center connection fails?**
    - Ensure configuration center service is running normally
    - Check network connection and firewall settings
    - System will automatically fall back to local configuration

2. **How to test hot configuration updates?**
    - After starting the application, directly modify values in remote configuration center
    - Wait for at least one `watchDuration` cycle
    - Observe application logs or behavior changes

3. **What configuration formats are supported?**
    - Supports mainstream formats including JSON, YAML, TOML
    - Specified by `configType` parameter

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/gone-io/goner/blob/main/LICENSE) file for details.