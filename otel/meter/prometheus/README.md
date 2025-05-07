<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/otel/meter/prometheus

## Overview

`goner/otel/meter/prometheus` is a Prometheus exporter component in the Gone framework that supports OpenTelemetry Metrics functionality. This module provides the capability to expose application metrics data in Prometheus format, enabling applications to easily integrate with the Prometheus monitoring system for metrics collection, storage, and visualization.

## Key Features

- Provides Prometheus format metrics reader
- Supports custom metric names and labels
- Supports various metric types (counters, gauges, histograms, etc.)
- Supports metric unit configuration
- Supports metric descriptions
- Integrates with Gone framework's lifecycle management

## Installation

```bash
# Install Prometheus metrics exporter
gonectl install goner/otel/meter/prometheus
```

## Example
> The following example demonstrates how to integrate OpenTelemetry with Prometheus in the Gone framework to implement application metrics monitoring and visualization. The project includes an HTTP service that uses Prometheus to scrape metrics data and Grafana for visualization.
> Full content: [Prometheus Metrics Monitoring](../../../examples/otel/meter/prometheus)

### Create Application Using gonectl
```bash
gonectl create -t otel/meter/prometheus prometheus-demo
cd prometheus-demo

# Start Prometheus
# make prometheus

# Start Application
# make run
```

### View Results

After the service is running, you can view the metrics data through the following methods:

1. Access the metrics endpoint: http://localhost:2112/metrics
2. Access Prometheus UI: http://localhost:9090
   - Enter metric name in the Graph interface
   - Click Execute button to view metrics data
3. Access Grafana interface: http://localhost:3000
   - Import the preset Dashboard
   - View metrics visualization panel

You can see complete metrics data, including:
- Basic application metrics
- Custom business metrics
- System resource usage
- Request processing statistics, etc.

Each metric contains detailed label information and help descriptions, facilitating data analysis and alert configuration.

## References

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Gone Framework Documentation](https://github.com/gone-io/gone)