[//]: # (desc: Using OpenTelemetry for Log Collection simple example)

<p>
    English &nbsp｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Using OpenTelemetry for Log Collection

This example demonstrates how to integrate OpenTelemetry's Log functionality in the Gone framework for application log collection.

## Project Setup Steps

### 1. Create Project and Install Dependencies

```bash
# Create project directory
mkdir simple-log
cd simple-log

# Initialize Go module
go mod init examples/otel/log/simple

# Install Gone framework's OpenTelemetry Log component
gonectl install goner/otel/log
```

### 2. Implement Main Program

Implement log recording in `main.go`:

```go
package main

import "github.com/gone-io/gone/v2"

func main() {
	gone.
		Loads(GoneModuleLoad).
		Run(func(logger gone.Logger) {
			logger.Infof("hello world")
			logger.Errorf("error info")
		})
}
```

### 3. Configure OpenTelemetry

Configure OpenTelemetry log collection in `config/default.yaml`:

```yaml
service:
  name: &serviceName "log-collect-example"

otel:
  service:
    name: *serviceName
  log:
    http:
      endpoint: localhost:4318
      insecure: true
  tracer:
    http:
      endpoint: localhost:4318
      insecure: true

log:
  otel:
    enable: true
    log-name: *serviceName
    only: false
```

## Run Service

```bash
# Run service
go run .
```

## View Results

After running the service, logs will be sent to the configured OpenTelemetry collector (endpoint: localhost:4318). You can view the collected logs in the OpenTelemetry collector's console or the configured storage backend (such as Elasticsearch, Loki, etc.).

The example logs include:
- An Info level log: "hello world"
- An Error level log: "error info"

These logs will be processed by the OpenTelemetry collector and include metadata information such as the service name (log-collect-example).