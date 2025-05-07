[//]: # (desc: OpenTelemetry Simple Tracing Example)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# OpenTelemetry Simple Tracing Example

This is a simple example of integrating OpenTelemetry tracing with the gone-io/goner framework.

## Features

This example demonstrates how to integrate OpenTelemetry tracing into a gone application, including:

- Injecting and using Tracer in components
- Creating Spans and adding events
- Passing context between function calls
- Configuring service name via environment variables

## Code Structure

- `main.go` - Application entry, sets service name and starts the app
- `your_component.go` - Example component, shows how to use OpenTelemetry Tracer
- `module.load.go` - Module loading configuration

## Code Example

### Main Program (main.go)

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"os"
)

func main() {
	// Set service name
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "simple demo")

	gone.
		Loads(GoneModuleLoad).
		Load(&YourComponent{}).
		Run(func(c *YourComponent) {
			// Call method in the component
			c.HandleRequest(context.Background())
		})
}
```

### Component Implementation (your_component.go)

```go
package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type YourComponent struct {
	gone.Flag
	tracer trace.Tracer `gone:"*,otel-tracer"` // Inject OpenTelemetry Tracer
}

func (c *YourComponent) HandleRequest(ctx context.Context) {
	tracer := otel.Tracer("otel-tracer")

	// Create a new Span
	ctx, span := tracer.Start(ctx, "handle-request")
	// Ensure the Span ends when the function returns
	defer span.End()

	// Add event
	span.AddEvent("Start handling request")

	// Business logic...

	// Record error (if any)
	// span.RecordError(err)
	// span.SetStatus(codes.Error, "Failed to handle request")

	// Set status to OK in normal case
	// span.SetStatus(codes.Ok, "")
}
```

## How to Run

1. Make sure Go is installed
2. Start OpenTelemetry Collector (or other compatible backend service)
3. Run the example:

```bash
go run .
```

## Dependencies

This example uses the following goner module:

- `github.com/gone-io/goner/otel/tracer` - OpenTelemetry tracing support

## More Information

For more information about OpenTelemetry, please refer to the [OpenTelemetry Official Documentation](https://opentelemetry.io/docs/).

For more information about the gone-io/goner framework, please refer to the [gone-io/goner Documentation](https://github.com/gone-io/gone).