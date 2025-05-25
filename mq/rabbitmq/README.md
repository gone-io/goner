<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/mq/rabbitmq Component and Gone RabbitMQ Integration

This package provides RabbitMQ integration for Gone applications, offering easy-to-use RabbitMQ client configuration and management based on the official RabbitMQ Go client library.

## Features

- Seamless integration with Gone's dependency injection system
- Support for multiple RabbitMQ client instances
- Provides producer and consumer interfaces
- Supports exchange and queue declaration
- Supports message routing and binding
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
    # Or use separate connection parameters:
    host: "localhost"                          # RabbitMQ host
    port: 5672                                 # RabbitMQ port
    username: "guest"                          # Authentication username
    password: "guest"                          # Authentication password
    vhost: "/"                                 # Virtual host
    producer:                                  # Producer configuration
      exchange: "my-exchange"                  # Exchange name
      exchangeType: "direct"                   # Exchange type (direct, fanout, topic, headers)
      durable: true                            # Whether the exchange is durable
      autoDelete: false                        # Whether the exchange is auto-deleted
      internal: false                          # Whether the exchange is internal
      noWait: false                            # Whether to wait for server confirmation
    consumer:                                  # Consumer configuration
      queue: "my-queue"                        # Queue name
      exchange: "my-exchange"                  # Bound exchange name
      routingKey: "my-routing-key"             # Bound routing key
      durable: true                            # Whether the queue is durable
      autoDelete: false                        # Whether the queue is auto-deleted
      exclusive: false                         # Whether the queue is exclusive
      noWait: false                            # Whether to wait for server confirmation
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
    producer goneRabbitMQ.IProducer `gone:"*"` // Default producer
}

func (p *Producer) SendMessage() {
    // Send message
    err := p.producer.Publish(
        "my-routing-key",
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType:  "text/plain",
            DeliveryMode: amqp.Persistent,
            Body:         []byte("Hello, RabbitMQ!"),
        },
    )
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
    consumer goneRabbitMQ.IConsumer `gone:"*"` // Default consumer
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
                log.Printf("Failed to acknowledge message: %v", err)
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

### Complete Example

Here is a complete producer and consumer example, demonstrating how to use them in the same application:

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
    producer goneRabbitMQ.IProducer `gone:"*"`
    consumer goneRabbitMQ.IConsumer `gone:"*"`
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
            err := app.producer.Publish(
                "",     // Routing key
                false, // mandatory
                false, // immediate
                amqp.Publishing{
                    ContentType:  "text/plain",
                    DeliveryMode: amqp.Persistent,
                    Timestamp:    time.Now(),
                    Body:         []byte(message),
                },
            )
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
                log.Printf("Failed to acknowledge message: %v", err)
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

## Environment Variable Configuration

You can also configure RabbitMQ through environment variables:

```bash
export GONE_RABBITMQ_DEFAULT='{"host":"127.0.0.1","port":5672,"username":"guest","password":"guest"}'
export GONE_RABBITMQ_DEFAULT_PRODUCER='{"exchange":"my-exchange","exchangeType":"direct","autoDelete":true,"durable":false,"noWait":false,"internal":false}'
export GONE_RABBITMQ_DEFAULT_CONSUMER='{"queue":"my-queue","exchange":"my-exchange","routingKey":"my-key","autoDelete":true,"durable":false,"noWait":false,"autoAck":true,"consumer":"my-consumer"}'
```

## Available Loaders

- `LoadConnection` - Loads only the RabbitMQ connection
- `LoadChannel` - Loads only the RabbitMQ channel (requires connection)
- `LoadProducer` - Loads the RabbitMQ producer (includes connection and channel)
- `LoadConsumer` - Loads the RabbitMQ consumer (includes connection and channel)
- `LoadAll` - Loads all RabbitMQ components