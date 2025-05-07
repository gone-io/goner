<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/otel/meter/prometheus/gin

## Overview

`goner/otel/meter/prometheus/gin` is a Gin integration component for the Prometheus exporter in the Gone framework that supports OpenTelemetry Metrics functionality. This module provides the capability to expose application metrics data in Prometheus format through Gin routing, enabling Gin-based applications to easily integrate with the Prometheus monitoring system for metrics collection, storage, and visualization.

## Key Features

- Provides Prometheus metrics exposure endpoint based on Gin
- Automatically registers `/metrics` route (configurable)
- Integrates with Gone framework's lifecycle management
- Easy-to-use interface with minimal configuration

## Installation

```bash
# Install the Prometheus metrics exporter's Gin integration component
gonectl install goner/otel/meter/prometheus/gin
```

## Basic Usage

### Loading the Module in Your Application

```go
func main() {
    gone.
		NewApp(gin.Load).
		Serve()
}
```

### Configuring the Metrics Endpoint

Add the following configuration to your configuration file to customize the metrics exposure path:

```yaml
otel:
  meter:
    prometheus:
      path: "/metrics"  # Prometheus scraping endpoint, defaults to /metrics
```

## Example

> The following example demonstrates how to integrate OpenTelemetry with Prometheus in the Gone framework to implement application metrics monitoring and visualization. The project includes a Gin-based HTTP service that exposes metrics data for Prometheus to scrape and uses Grafana for visualization.
> Full content: [Prometheus Metrics Monitoring](../../../../examples/otel/meter/prometheus)

### Creating a Project and Installing Dependencies

```bash
# Create project directory
mkdir prometheus-demo
cd prometheus-demo

# Initialize Go module
go mod init examples/prometheus-demo

# Install the Prometheus Gin integration component
gonectl install goner/otel/meter/prometheus/gin
```

### Implementing Custom Metrics

Create a controller file to implement an API access counter:

```go
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type ctr struct {
	gone.Flag
	r g.IRoutes `gone:"*"`
}

func (c *ctr) Mount() (err g.MountError) {
	var meter = otel.Meter("my-service-meter")
	apiCounter, err := meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls"),
		metric.WithUnit("{count}"),
	)
	if err != nil {
		return gone.ToErrorWithMsg(err, "Failed to create api.counter")
	}

	c.r.GET("/hello", func(ctx *gin.Context) string {
		apiCounter.Add(ctx, 1)
		return "hello, world"
	})
	return
}
```

### Viewing Results

After the service is running, you can view the metrics data through:

1. Access the metrics endpoint: http://localhost:2112/metrics
2. Access the Prometheus UI: http://localhost:9090
   - Enter the metric name (e.g., `api_counter`) in the Graph interface
   - Click the Execute button to view metric data
3. Access the Grafana interface: http://localhost:3000
   - Import the preset Dashboard
   - View the metrics visualization panel

## How It Works

The `goner/otel/meter/prometheus/gin` module works in the following way:

1. During application startup, it automatically registers a Gin route handler to expose metrics data in Prometheus format
2. When the Prometheus server accesses the configured endpoint (default `/metrics`), the module collects all current application metrics and returns them in Prometheus format
3. Applications can use the OpenTelemetry API to create and update various types of metrics (counters, gauges, histograms, etc.)

## References

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [Gin Web Framework](https://github.com/gin-gonic/gin)