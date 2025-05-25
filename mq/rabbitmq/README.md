# goner/mq/rabbitmq Component and Gone RabbitMQ Integration

This package provides RabbitMQ integration functionality for Gone applications, offering simple and easy-to-use RabbitMQ client configuration and management based on the official RabbitMQ Go client library.

## Features

- Easy integration with Gone's dependency injection system
- Support for multiple RabbitMQ client instances
- Provides producers and consumers
- Support for exchange and queue declarations
- Support for message routing and binding
- Automatic resource management and cleanup
- Comprehensive configuration options

## Installation

```bash
gonectl install goner/mq/rabbitmq
gonectl install goner/viper # Optional, for loading configuration files
```

## Configuration

Create a `config/default.yaml` file in your project's configuration directory and add the following RabbitMQ configuration:

```yaml
rabbitmq:
  default:
    url: "amqp://guest:guest@localhost:5672/"  # RabbitMQ connection URL
    # OR use individual connection parameters:
    host: "localhost"                          # RabbitMQ host
    port: 5672                                 # RabbitMQ port
    username: "guest"                          # Username for authentication
    password: "guest"                          # Password for authentication
    vhost: "/"                                 # Virtual host
    producer:                                  # Producer configuration
      exchange: "my-exchange"                  # Exchange name
      exchangeType: "direct"                   # Exchange type (direct, fanout, topic, headers)
      routingKey: "my-routing-key"             # Default routing key
      durable: true                            # Whether exchange is durable
      autoDelete: false                        # Whether exchange auto-deletes
    consumer:                                  # Consumer configuration
      queue: "my-queue"                        # Queue name
      exchange: "my-exchange"                  # Exchange name to bind to
      routingKey: "my-routing-key"             # Routing key for binding
      durable: true                            # Whether queue is durable
      autoDelete: false                        # Whether queue auto-deletes
      exclusive: false                         # Whether queue is exclusive
      autoAck: false                           # Whether to auto-acknowledge messages
      consumer: "my-consumer"                  # Consumer tag
```

### Multiple Client Configuration

You can configure multiple RabbitMQ client instances:

```yaml
rabbitmq:
  default:
    url: "amqp://guest:guest@localhost:5672/"
  cluster1:
    host: "rabbitmq1.example.com"
    port: 5672
    username: "user1"
    password: "pass1"
    vhost: "/production"
  cluster2:
    host: "rabbitmq2.example.com"
    port: 5672
    username: "user2"
    password: "pass2"
    vhost: "/testing"
```

## Usage

### Basic Usage

#### Producer

```go
package main

import (
    "github.com/gone-io/gone/v2"
    goneRabbitMQ "github.com/gone-io/goner/mq/rabbitmq"
    amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
    gone.Flag
    producer *goneRabbitMQ.Producer `gone:"*"` // Default producer
}

func (p *Producer) SendMessage() {
    // Send message
    err := p.producer.Publish("my-routing-key", []byte("Hello, RabbitMQ!"), amqp.Table{
        "custom-header": "value",
    })
    if err != nil {
        // Handle error
        return
    }
    
    println("Message sent successfully")
}

func main() {
    gone.NewApp(goneRabbitMQ.LoadProducer).Run(func(p *Producer) {
        p.SendMessage()
    })
}
```

#### Consumer

```go
package main

import (
    "context"
    "log"
    "github.com/gone-io/gone/v2"
    goneRabbitMQ "github.com/gone-io/goner/mq/rabbitmq"
)

type Consumer struct {
    gone.Flag
    consumer *goneRabbitMQ.Consumer `gone:"*"` // Default consumer
}

func (c *Consumer) ConsumeMessages(ctx context.Context) {
    // Start consuming messages
    msgs, err := c.consumer.Consume()
    if err != nil {
        log.Fatal("Failed to start consuming:", err)
    }
    
    for {
        select {
        case msg := <-msgs:
            // Process message
            log.Printf("Received message: %s", string(msg.Body))
            
            // Acknowledge message (if autoAck is false)
            if err := msg.Ack(false); err != nil {
                log.Printf("Failed to ack message: %v", err)
            }
        case <-ctx.Done():
            return
        }
    }
}

func main() {
    gone.NewApp(goneRabbitMQ.LoadConsumer).Run(func(c *Consumer) {
        ctx := context.Background()
        c.ConsumeMessages(ctx)
    })
}
```

### Multiple Client Usage

To use multiple RabbitMQ clients, specify the client name in the Gone tag:

```go
type MultiClientApp struct {
    gone.Flag
    defaultProducer *goneRabbitMQ.Producer `gone:"*"`         // Default client
    cluster1Producer *goneRabbitMQ.Producer `gone:"cluster1"`  // cluster1 client
    cluster2Consumer *goneRabbitMQ.Consumer `gone:"cluster2"`  // cluster2 client
}
```

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "time"
    
    "github.com/gone-io/gone/v2"
    goneRabbitMQ "github.com/gone-io/goner/mq/rabbitmq"
    amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQApp struct {
    gone.Flag
    producer *goneRabbitMQ.Producer `gone:"*"`
    consumer *goneRabbitMQ.Consumer `gone:"*"`
}

func (app *RabbitMQApp) Run() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Send messages
    go app.sendMessages(ctx)
    
    // Consume messages
    go app.consumeMessages(ctx)
    
    // Wait for interrupt signal
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    <-signals
}

func (app *RabbitMQApp) sendMessages(ctx context.Context) {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    i := 0
    for {
        select {
        case <-ticker.C:
            message := fmt.Sprintf("Message %d", i)
            err := app.producer.Publish("", []byte(message), amqp.Table{
                "timestamp": time.Now().Unix(),
            })
            if err != nil {
                log.Printf("Failed to send message: %v", err)
            } else {
                log.Printf("Sent: %s", message)
            }
            i++
        case <-ctx.Done():
            return
        }
    }
}

func (app *RabbitMQApp) consumeMessages(ctx context.Context) {
    msgs, err := app.consumer.Consume()
    if err != nil {
        log.Fatal("Failed to start consuming:", err)
    }
    
    for {
        select {
        case msg := <-msgs:
            log.Printf("Received: %s", string(msg.Body))
            
            // Acknowledge message
            if err := msg.Ack(false); err != nil {
                log.Printf("Failed to ack message: %v", err)
            }
        case <-ctx.Done():
            return
        }
    }
}

func main() {
    gone.NewApp(
        goneRabbitMQ.LoadProducer,
        goneRabbitMQ.LoadConsumer,
    ).Run(func(app *RabbitMQApp) {
        app.Run()
    })
}
```

## Advanced Usage

### Environment Variable Configuration

You can also configure RabbitMQ using environment variables:

```bash
export GONE_RABBITMQ_DEFAULT='{"url":"amqp://guest:guest@localhost:5672/"}'
export GONE_RABBITMQ_DEFAULT_PRODUCER='{"exchange":"my-exchange","routingKey":"my-key"}'
export GONE_RABBITMQ_DEFAULT_CONSUMER='{"queue":"my-queue","exchange":"my-exchange"}'
```

### Custom Configuration

The component supports all RabbitMQ configuration options provided by the `github.com/rabbitmq/amqp091-go` package, including:

- Connection configuration (URL, host, port, credentials, virtual host)
- Producer configuration (exchange declaration, routing, durability)
- Consumer configuration (queue declaration, binding, acknowledgment)
- Exchange and queue properties (durable, auto-delete, exclusive)
- Message properties and headers

For detailed configuration options, please refer to the [RabbitMQ Go Client Documentation](https://github.com/rabbitmq/amqp091-go).

## Available Loaders

- `LoadConnection` - Load RabbitMQ connection only
- `LoadChannel` - Load RabbitMQ channel only (requires connection)
- `LoadProducer` - Load RabbitMQ producer (includes connection and channel)
- `LoadConsumer` - Load RabbitMQ consumer (includes connection and channel)
- `LoadAll` - Load all RabbitMQ components