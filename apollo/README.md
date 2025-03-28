# Gone Apollo Component

## Introduction

The Gone Apollo component is an integration of the [Apollo](https://www.apolloconfig.com/) configuration center with the Gone framework, providing dynamic configuration fetching and real-time updates. Apollo is an open-source distributed configuration center developed by Ctrip, which can centrally manage configurations for different environments and clusters. After configuration changes, they can be pushed to applications in real-time, with standardized permissions and governance features.

## Quick Start

### 1. Load Apollo Configuration Component

```go
package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/apollo"
)

func main() {
	gone.
		Loads(
			apollo.Load, // Load Apollo configuration component
			// Other components...
		).
		// Or Serve()
		Run()
}
```

### 2. Configure Apollo Connection Information

Add the following configuration in your project's configuration file (e.g., `config/default.yaml`):

```yaml
apollo.appId: YourAppId           # Apollo application ID
apollo.cluster: default           # Cluster name, default is 'default'
apollo.ip: http://apollo-server:8080  # Apollo configuration center address
apollo.namespace: application     # Namespace, default is 'application'
apollo.secret: YourSecretKey      # Access key (if access key verification is enabled)
apollo.isBackupConfig: true       # Whether to enable backup configuration
apollo.watch: true                # Whether to listen for configuration changes
apollo.useLocalConfIfKeyNotExist: true  # If a key doesn't exist in Apollo config, whether to use the value from local config file
```

### 3. Use Configuration

Inject configuration in Gone components:

```go
type YourComponent struct {
	gone.Flag
	
	// Method 1: Directly inject configuration value
	DbUrl string `gone:"config,database.url"`
	
	// Method 2: Get configuration through Configure interface
	configure gone.Configure `gone:"*"`
}

func (c *YourComponent) AfterProp() {
	// Method 2: Dynamically get configuration
	var port int
	err := c.configure.Get("server.port", &port, "8080")
	if err != nil {
		// Handle error
	}
}
```

## Dynamic Configuration Updates

When `apollo.watch` is set to `true`, the Apollo component will listen for configuration changes and automatically update registered configuration items.
**Note**: For fields that need dynamic updates, **pointer types must be used** to take effect.

To make configuration items support dynamic updates, they need to be registered with the change listener when getting the configuration:

```go
type YourComponent struct {
	gone.Flag
	
	// These configuration items will support dynamic updates
	ServerPort *int    `gone:"config,server.port"`
	DbUrl      *string `gone:"config,database.url"`
}

// After configuration changes, ServerPort and DbUrl values will be automatically updated
```

## Configuration Items

| Configuration Item | Description | Default Value |
| --- | --- | --- |
| apollo.appId | Apollo application ID, must match the application ID in Apollo configuration center | - |
| apollo.cluster | Cluster name | default |
| apollo.ip | Apollo configuration center address | - |
| apollo.namespace | Namespace | application |
| apollo.secret | Access key for client authentication | - |
| apollo.isBackupConfig | Whether to enable backup configuration, enabling will save config locally | true |
| apollo.watch | Whether to listen for configuration changes, enabling will auto-update on changes | false |
|apollo.useLocalConfIfKeyNotExist|If a key doesn't exist in Apollo config, whether to use the value from local config file|true|

## Advanced Usage

### Multiple Namespaces Support

Apollo supports multiple namespaces, default is `application`. To use multiple namespaces, specify in configuration:

```yaml
apollo.namespace: application,common,custom
```

### Local Cache Configuration

When `apollo.isBackupConfig` is set to `true`, Apollo client will cache configurations locally. When Apollo service is unavailable, local cached configurations will be used.

## Notes

1. Ensure Apollo configuration center is properly deployed and accessible
2. Configuration type conversion is handled by Gone framework, supports basic types (string, int, bool, etc.)
3. For complex types (structs, arrays, etc.), Apollo client will try to parse configuration values as JSON
4. Configuration change listening requires setting `apollo.watch: true`

## References

- [Apollo Official Documentation](https://www.apolloconfig.com/)
- [Gone Framework Documentation](https://github.com/gone-io/gone)