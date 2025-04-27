[//]: # (desc: Etcd Configuration Center Example)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone Framework Etcd Configuration Center Example

- [Gone Framework Etcd Configuration Center Example](#gone-framework-etcd-configuration-center-example)
  - [Overview](#overview)
  - [Environment Setup](#environment-setup)
    - [Starting the Environment](#starting-the-environment)
  - [Configuration File Structure](#configuration-file-structure)
    - [Local Configuration File](#local-configuration-file)
    - [Configuration Files in Etcd](#configuration-files-in-etcd)
  - [Code Implementation](#code-implementation)
  - [Running the Example](#running-the-example)
    - [1. Import Configuration to Etcd](#1-import-configuration-to-etcd)
    - [2. Run the Example Program](#2-run-the-example-program)
    - [3. Test Dynamic Configuration Updates](#3-test-dynamic-configuration-updates)
  - [Configuration Priority](#configuration-priority)
  - [Summary](#summary)


This example demonstrates how to integrate the Gone framework with Etcd configuration center to achieve centralized configuration management and dynamic updates.

## Overview

This example demonstrates the following features:

- Using Gone framework to read configurations from Etcd configuration center
- Automatic configuration monitoring and dynamic updates
- Structured configuration injection
- Mixed use of local and remote configurations

## Environment Setup

This example uses Docker Compose to start Etcd service and Etcd management interface (etcdkeeper):

```yaml
services:
  Etcd:
    image: 'bitnami/etcd:latest'
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379
    ports:
      - "2379:2379"
      - "2380:2380"
  etcdKeeper:
    image: evildecay/etcdkeeper
    environment:
      HOST: "0.0.0.0"
    ports:
      - "12000:8080"
    depends_on:
      - Etcd
```

### Starting the Environment

```bash
docker-compose up -d
```

After startup, you can access the etcdkeeper management interface at http://localhost:12000.

## Configuration File Structure

### Local Configuration File

`config/default.yaml` contains Viper remote configuration settings and local configuration:

```yaml
viper.remote:
  type: yaml
  watch: true
  watchDuration: 5s
  useLocalConfIfKeyNotExist: true
  providers:
    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path:  /config/application.yaml
      keyring:

    - provider: etcd3
      configType: yaml
      endpoint: http://localhost:2379
      path:  /config/database.yaml
      keyring:


key:
  not-existed-in-etcd: 1000
```

Configuration explanation:
- `watch: true` - Enable configuration change monitoring
- `watchDuration: 5s` - Monitoring interval is 5 seconds
- `useLocalConfIfKeyNotExist: true` - Use local configuration when a key doesn't exist in remote configuration
- Configured two Etcd configuration sources: `/config/application.yaml` and `/config/database.yaml`

### Configuration Files in Etcd

Configuration files from the `etcd-config-files` directory need to be imported into Etcd:

**application.yaml**:
```yaml
# /config/application.yaml

server.name: config-demo
server.port: 9090
```

**database.yaml**:
```yaml
# /config/database.yaml

database:
  username: config-demo
  password: config-demo-password
```

## Code Implementation

The main program `main.go` demonstrates how to inject configurations using the Gone framework:

```go
package main

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/viper/remote"
	"time"
)

type Database struct {
	UserName string `mapstructure:"username"`
	Pass     string `mapstructure:"password"`
}

func main() {
	gone.
		NewApp(remote.Load).
		Run(func(params struct {
			serverName string `gone:"config,server.name"`
			serverPort int    `gone:"config,server.port"`

			dbUserName string `gone:"config,database.username"`
			dbUserPass string `gone:"config,database.password"`

			database *Database `gone:"config,database"`

			key string `gone:"config,key.not-existed-in-etcd"`
		}) {
			fmt.Printf("serverName=%s, serverPort=%d, dbUserName=%s, dbUserPass=%s, key=%s\n", params.serverName, params.serverPort, params.dbUserName, params.dbUserPass, params.key)

			for i := 0; i < 10; i++ {
				fmt.Printf("database: %#+v\n", *params.database)
				time.Sleep(10 * time.Second)
			}
		})
}
```

Code explanation:

1. Load remote configuration component through `remote.Load`
2. Use `gone:"config,xxx"` tags to inject configuration items:
   - Basic type configurations: `serverName`, `serverPort`, etc.
   - Struct configuration: `database`
   - Local configuration: `key.not-existed-in-etcd` (exists only in local configuration)
3. The program prints database configuration every 10 seconds to demonstrate dynamic configuration updates

## Running the Example

### 1. Import Configuration to Etcd

Use etcdkeeper (http://localhost:12000) to import configuration files from the `etcd-config-files` directory into Etcd:

- Create key `/config/application.yaml` with the content of application.yaml
- Create key `/config/database.yaml` with the content of database.yaml

### 2. Run the Example Program

```bash
go run main.go
```

### 3. Test Dynamic Configuration Updates

1. After running the program, modify the configuration in `/config/database.yaml` through etcdkeeper
2. Observe the program output, configuration will automatically update after about 5 seconds

## Configuration Priority

1. Remote configuration (Etcd) has higher priority than local configuration
2. When a key doesn't exist in remote configuration, local configuration will be used (controlled by `useLocalConfIfKeyNotExist: true`)

## Summary

This example demonstrates how the Gone framework integrates with Etcd configuration center to achieve centralized configuration management and dynamic updates. Through simple configuration and minimal code, powerful configuration management functionality can be implemented, providing a flexible configuration solution for microservice architecture.