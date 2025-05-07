<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/otel/tracer/zipkin

## Overview

`goner/otel/tracer/zipkin` is an OpenTelemetry exporter component in the Gone framework that supports the Zipkin tracing system. This module provides functionality to export OpenTelemetry trace data to the Zipkin tracing system, enabling applications to integrate with the existing Zipkin ecosystem.

## Key Features

- Provides Zipkin format OpenTelemetry trace exporter
- Supports custom HTTP headers
- Integrates with Gone framework's lifecycle management
- Compatible with Zipkin tracing system's data format and API

## Installation

```bash
# Install Zipkin trace exporter
gonectl install goner/otel/tracer/zipkin
```

## Configuration

To use the Zipkin trace exporter in your application, add the following configuration to your Gone framework configuration file:

```yaml
otel:
  service:
    name: "your-service-name"  # Set service name
  tracer:
    zipkin:
      url: "http://your-zipkin-endpoint/api/v2/spans"  # Zipkin receiver endpoint
      headers:  # Optional, custom HTTP headers
        Authorization: "Bearer your-token"
```

### Configuration Options

| Option | Type | Description |
| --- | --- | --- |
| `url` | string | Complete URL of the Zipkin receiver endpoint |
| `headers` | map | Custom HTTP headers |



## References

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Zipkin Documentation](https://zipkin.io/)
- [Gone Framework Documentation](https://github.com/gone-io/gone)