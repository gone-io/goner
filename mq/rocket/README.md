<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/mq/rocket Component and Gone RocketMQ Integration

This package provides RocketMQ integration functionality for Gone applications, offering simple and easy-to-use RocketMQ client configuration and management based on the Apache RocketMQ Go client library.

## Features

- Easy integration with Gone's dependency injection system
- Support for multiple RocketMQ client instances
- Provides producers and simple consumers
- Support for transactional messages
- Support for SQL and Tag filtering expressions
- Automatic resource management and cleanup
- Comprehensive configuration options

## Installation

```bash
gonectl install goner/mq/rocket
gonectl install goner/viper # Optional, for loading configuration files
```

## Configuration

Create a `config/default.yaml` file in your project's configuration directory and add the following RocketMQ configuration:

```yaml
rocketmq:
  default:
    Endpoint: "127.0.0.1:8081"          # RocketMQ broker server address
    ConsumerGroup: "my-consumer-group"    # Consumer group ID
    NameSpace: ""                        # Namespace (optional)
    Credentials:                          # Authentication information (optional)
      accessKey: ""
      accessSecret: ""
      securityToken: ""
    consumer:                             # Consumer configuration
      awaitDuration: 10000000000          # Wait duration (nanoseconds)
      expressions:                        # Subscription expressions
        - topic: "my-topic"
          type: "tag"                      # Filter type: tag or sql
          expression: "*"                  # Filter expression
    producer:                             # Producer configuration
      maxAttempts: 3                      # Maximum retry attempts
      topics: ["my-topic"]               # Predefined topic list
```

### Multiple Client Configuration

You can configure multiple RocketMQ client instances:

```yaml
rocketmq:
  default:
    Endpoint: "127.0.0.1:8081"
    ConsumerGroup: "default-group"
  cluster1:
    Endpoint: "rocketmq1.example.com:8081"
    ConsumerGroup: "cluster1-group"
    NameSpace: "production"
  cluster2:
    Endpoint: "rocketmq2.example.com:8081"
    ConsumerGroup: "cluster2-group"
    NameSpace: "testing"
```

## Usage

### Basic Usage

#### Producer

```go
package main

import (
    "context"
    mq "github.com/apache/rocketmq-clients/golang/v5"
    "github.com/gone-io/gone/v2"
    goneRocket "github.com/gone-io/goner/mq/rocket"
)

type Producer struct {
    gone.Flag
    producer mq.Producer `gone:"*"` // Default producer
}

func (p *Producer) SendMessage() {
    msg := &mq.Message{
        Topic: "my-topic",
        Body:  []byte("Hello, RocketMQ!"),
    }
    
    // Set message properties
    msg.SetKeys("key1", "key2")
    msg.SetTag("important")
    
    // Send message
    resp, err := p.producer.Send(context.TODO(), msg)
    if err != nil {
        // Handle error
        return
    }
    
    // Handle send result
    for _, result := range resp {
        println("Message sent, MessageID:", result.GetMessageId())
    }
}

func main() {
    gone.NewApp(goneRocket.LoadProducer).Run(func(p *Producer) {
        p.SendMessage()
    })
}
```

#### Simple Consumer

```go
type Consumer struct {
    gone.Flag
    consumer mq.SimpleConsumer `gone:"*"` // Default simple consumer
}

func (c *Consumer) ConsumeMessages() {
    ctx := context.Background()
    maxMessageNum := int32(16)
    invisibleDuration := time.Second * 20
    
    for {
        // Receive messages
        messages, err := c.consumer.Receive(ctx, maxMessageNum, invisibleDuration)
        if err != nil {
            println("Failed to receive messages:", err.Error())
            continue
        }
        
        // Process messages
        for _, message := range messages {
            println("Received message:", string(message.GetBody()))
            
            // Acknowledge message
            if err := c.consumer.Ack(ctx, message); err != nil {
                println("Failed to ack message:", err.Error())
            }
        }
        
        time.Sleep(time.Second)
    }
}
```

### Multiple Clients Example

```go
type MultiRocketMQClient struct {
    gone.Flag
    defaultProducer mq.Producer       `gone:"*"`           // Default producer
    cluster1Producer mq.Producer      `gone:"*,cluster1"`  // cluster1 producer
    cluster2Consumer mq.SimpleConsumer `gone:"*,cluster2"`  // cluster2 consumer
}

func (m *MultiRocketMQClient) UseMultipleClients() {
    ctx := context.Background()
    
    // Send message using default cluster
    msg1 := &mq.Message{
        Topic: "topic1",
        Body:  []byte("Message to default cluster"),
    }
    m.defaultProducer.Send(ctx, msg1)
    
    // Send message using cluster1
    msg2 := &mq.Message{
        Topic: "topic2",
        Body:  []byte("Message to cluster1"),
    }
    m.cluster1Producer.Send(ctx, msg2)
    
    // Consume messages from cluster2
    messages, _ := m.cluster2Consumer.Receive(ctx, 10, time.Second*20)
    for _, message := range messages {
        println("Received from cluster2:", string(message.GetBody()))
        m.cluster2Consumer.Ack(ctx, message)
    }
}
```

### Filtering Expression Examples

#### Tag Filtering

```yaml
rocketmq:
  default:
    consumer:
      expressions:
        - topic: "order-topic"
          type: "tag"
          expression: "important || urgent"  # Receive messages with important or urgent tags
```

#### SQL Filtering

```yaml
rocketmq:
  default:
    consumer:
      expressions:
        - topic: "user-topic"
          type: "sql"
          expression: "age >= 18 AND region = 'US'"  # SQL expression filtering
```

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "time"
    mq "github.com/apache/rocketmq-clients/golang/v5"
    "github.com/gone-io/gone/v2"
    goneRocket "github.com/gone-io/goner/mq/rocket"
)

type RocketMQApp struct {
    gone.Flag
    producer mq.Producer       `gone:"*"`
    consumer mq.SimpleConsumer `gone:"*"`
}

func (app *RocketMQApp) Run() {
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

func (app *RocketMQApp) sendMessages(ctx context.Context) {
    for i := 0; i < 10; i++ {
        msg := &mq.Message{
            Topic: "test-topic",
            Body:  []byte(fmt.Sprintf("Message %d", i)),
        }
        
        // Set message properties
        msg.SetKeys(fmt.Sprintf("key-%d", i))
        msg.SetTag("test")
        
        resp, err := app.producer.Send(ctx, msg)
        if err != nil {
            println("Failed to send message:", err.Error())
        } else {
            for _, result := range resp {
                println("Message sent, ID:", result.GetMessageId())
            }
        }
        
        time.Sleep(time.Second)
    }
}

func (app *RocketMQApp) consumeMessages(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            messages, err := app.consumer.Receive(ctx, 16, time.Second*20)
            if err != nil {
                println("Failed to receive messages:", err.Error())
                continue
            }
            
            for _, message := range messages {
                println("Received:", string(message.GetBody()))
                
                // Acknowledge message
                if err := app.consumer.Ack(ctx, message); err != nil {
                    println("Failed to ack message:", err.Error())
                }
            }
        }
    }
}

func main() {
    gone.NewApp(
        goneRocket.LoadProducer,
        goneRocket.LoadConsumer,
    ).Run(func(app *RocketMQApp) {
        app.Run()
    })
}
```

## Advanced Usage

### Environment Variable Configuration

You can also configure RocketMQ using environment variables:

```bash
export GONE_ROCKETMQ_DEFAULT='{"Endpoint":"127.0.0.1:8081","ConsumerGroup":"my-group"}'
export GONE_ROCKETMQ_DEFAULT_CONSUMER='{"awaitDuration":10000000000,"expressions":[{"topic":"my-topic"}]}'
export GONE_ROCKETMQ_DEFAULT_PRODUCER='{"maxAttempts":3,"topics":["my-topic"]}'
```

### Transaction Message Support

The component supports transaction messages. You can handle transaction checks by registering a `TransactionChecker`:

```go
type MyTransactionChecker struct {
    gone.Flag
}

func (c *MyTransactionChecker) Check(msg *mq.MessageView) mq.TransactionResolution {
    // Implement transaction status check logic
    return mq.COMMIT
}

// Register transaction checker
func (c *MyTransactionChecker) GetGoneId() string {
    return "rocketmq.default" // Corresponds to client name in configuration
}
```

### Custom Configuration

The component supports all RocketMQ configuration options provided by the `github.com/apache/rocketmq-clients/golang/v5` package, including:

- Producer configuration (retry attempts, topic list, etc.)
- Consumer configuration (wait duration, filter expressions, etc.)
- Network configuration (endpoints, timeouts, etc.)
- Security configuration (authentication information, etc.)
- Namespace configuration

For detailed configuration options, please refer to the [RocketMQ Go Client Documentation](https://github.com/apache/rocketmq-clients/tree/master/golang).

## Available Loaders

- `LoadConsumer` - Load RocketMQ simple consumer
- `LoadProducer` - Load RocketMQ producer

## Important Notes

1. All client instances will automatically clean up resources when the application shuts down
2. Consumers require the `ConsumerGroup` parameter to be configured
3. It is recommended to configure appropriate retry and timeout parameters in production environments
4. When using SQL filtering expressions, ensure that the RocketMQ server supports this feature
5. Transaction messages require implementing the `TransactionChecker` interface to handle transaction checks
6. The consumer's `awaitDuration` parameter is in nanoseconds
