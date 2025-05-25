<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/mq/kafka Component and Gone Kafka Integration

This package provides Kafka integration functionality for Gone applications, offering simple and easy-to-use Kafka client configuration and management based on the IBM Sarama library.

## Features

- Easy integration with Gone's dependency injection system
- Support for multiple Kafka client instances
- Provides synchronous and asynchronous producers
- Support for consumers and consumer groups
- Automatic resource management and cleanup
- Comprehensive configuration options

## Installation

```bash
gonectl install goner/mq/kafka
gonectl install goner/viper # Optional, for loading configuration files
```

## Configuration

Create a `config/default.yaml` file in your project's configuration directory and add the following Kafka configuration:

```yaml
kafka:
  default:
    addrs: ["127.0.0.1:9092"]     # List of Kafka server addresses
    groupID: "my-consumer-group"   # Consumer group ID
    Producer:
      Return:
        Successes: true            # Whether to return success messages
        Errors: true               # Whether to return error messages
      RequiredAcks: 1              # Acknowledgment level
      Retry:
        Max: 3                     # Maximum retry attempts
    Consumer:
      Offsets:
        AutoCommit:
          Enable: true             # Whether to auto-commit offsets
          Interval: 1000           # Auto-commit interval (milliseconds)
        Initial: -2                # Initial offset (-2: earliest, -1: latest)
      Group:
        Rebalance:
          Strategy: "range"        # Rebalance strategy
```

### Multiple Client Configuration

You can configure multiple Kafka client instances:

```yaml
kafka:
  default:
    addrs: ["127.0.0.1:9092"]
    groupID: "default-group"
  cluster1:
    addrs: ["kafka1.example.com:9092", "kafka2.example.com:9092"]
    groupID: "cluster1-group"
  cluster2:
    addrs: ["kafka3.example.com:9092"]
    groupID: "cluster2-group"
```

## Usage

### Basic Usage

#### Synchronous Producer

```go
package main

import (
    "github.com/IBM/sarama"
    "github.com/gone-io/gone/v2"
    goneKafka "github.com/gone-io/goner/mq/kafka"
)

type Producer struct {
    gone.Flag
    syncProducer sarama.SyncProducer `gone:"*"` // Default synchronous producer
}

func (p *Producer) SendMessage() {
    msg := &sarama.ProducerMessage{
        Topic: "my-topic",
        Value: sarama.StringEncoder("Hello, Kafka!"),
    }
    
    partition, offset, err := p.syncProducer.SendMessage(msg)
    if err != nil {
        // Handle error
        return
    }
    // Message sent successfully
    println("Message sent to partition", partition, "at offset", offset)
}

func main() {
    gone.NewApp(goneKafka.LoaderSyncProducer).Run(func(p *Producer) {
        p.SendMessage()
    })
}
```

#### Asynchronous Producer

```go
type AsyncProducer struct {
    gone.Flag
    asyncProducer sarama.AsyncProducer `gone:"*"` // Default asynchronous producer
}

func (p *AsyncProducer) SendMessage() {
    msg := &sarama.ProducerMessage{
        Topic: "my-topic",
        Value: sarama.StringEncoder("Hello, Async Kafka!"),
    }

    // Send message to input channel
    p.asyncProducer.Input() <- msg

    // Listen to success and error channels
    go func() {
        for {
            select {
            case success := <-p.asyncProducer.Successes():
                println("Message sent successfully:", success.Offset)
            case err := <-p.asyncProducer.Errors():
                println("Failed to send message:", err.Err)
            }
        }
    }()
}
```

#### Consumer

```go
type Consumer struct {
    gone.Flag
    consumer sarama.Consumer `gone:"*"` // Default consumer
}

func (c *Consumer) ConsumeMessages() {
    partitionConsumer, err := c.consumer.ConsumePartition("my-topic", 0, sarama.OffsetNewest)
    if err != nil {
        // Handle error
        return
    }
    defer partitionConsumer.Close()
    
    for {
        select {
        case message := <-partitionConsumer.Messages():
            println("Received message:", string(message.Value))
        case err := <-partitionConsumer.Errors():
            println("Error:", err.Err)
        }
    }
}
```

#### Consumer Group

```go
type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
    return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for {
        select {
        case message := <-claim.Messages():
            println("Received message:", string(message.Value))
            session.MarkMessage(message, "")
        case <-session.Context().Done():
            return nil
        }
    }
}

type ConsumerGroup struct {
    gone.Flag
    consumerGroup sarama.ConsumerGroup `gone:"*"` // Default consumer group
}

func (c *ConsumerGroup) ConsumeMessages(ctx context.Context) {
    handler := &ConsumerGroupHandler{}
    
    for {
        if err := c.consumerGroup.Consume(ctx, []string{"my-topic"}, handler); err != nil {
            // Handle error
            return
        }
        
        if ctx.Err() != nil {
            return
        }
    }
}
```

### Multiple Clients Example

```go
type MultiKafkaClient struct {
    gone.Flag
    defaultProducer sarama.SyncProducer `gone:"*"`           // Default producer
    cluster1Producer sarama.SyncProducer `gone:"*,cluster1"`  // cluster1 producer
    cluster2Consumer sarama.Consumer     `gone:"*,cluster2"`  // cluster2 consumer
}

func (m *MultiKafkaClient) UseMultipleClients() {
    // Send message using default cluster
    msg1 := &sarama.ProducerMessage{
        Topic: "topic1",
        Value: sarama.StringEncoder("Message to default cluster"),
    }
    m.defaultProducer.SendMessage(msg1)
    
    // Send message using cluster1
    msg2 := &sarama.ProducerMessage{
        Topic: "topic2",
        Value: sarama.StringEncoder("Message to cluster1"),
    }
    m.cluster1Producer.SendMessage(msg2)
    
    // Consume messages from cluster2
    partitionConsumer, _ := m.cluster2Consumer.ConsumePartition("topic3", 0, sarama.OffsetNewest)
    defer partitionConsumer.Close()
    
    for message := range partitionConsumer.Messages() {
        println("Received from cluster2:", string(message.Value))
    }
}
```

### Complete Example

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "github.com/IBM/sarama"
    "github.com/gone-io/gone/v2"
    goneKafka "github.com/gone-io/goner/mq/kafka"
)

type KafkaApp struct {
    gone.Flag
    syncProducer  sarama.SyncProducer  `gone:"*"`
    asyncProducer sarama.AsyncProducer `gone:"*"`
    consumer      sarama.Consumer      `gone:"*"`
    consumerGroup sarama.ConsumerGroup `gone:"*"`
}

func (app *KafkaApp) Run() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Send synchronous messages
    go app.sendSyncMessages()
    
    // Send asynchronous messages
    go app.sendAsyncMessages()
    
    // Consume messages
    go app.consumeMessages(ctx)
    
    // Wait for interrupt signal
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    <-signals
}

func (app *KafkaApp) sendSyncMessages() {
    for i := 0; i < 10; i++ {
        msg := &sarama.ProducerMessage{
            Topic: "test-topic",
            Value: sarama.StringEncoder(fmt.Sprintf("Sync message %d", i)),
        }
        partition, offset, err := app.syncProducer.SendMessage(msg)
        if err != nil {
            println("Failed to send sync message:", err.Error())
        } else {
            println("Sync message sent to partition", partition, "at offset", offset)
        }
    }
}

func (app *KafkaApp) sendAsyncMessages() {
    for i := 0; i < 10; i++ {
        msg := &sarama.ProducerMessage{
            Topic: "test-topic",
            Value: sarama.StringEncoder(fmt.Sprintf("Async message %d", i)),
        }
        app.asyncProducer.Input() <- msg
    }
}

func (app *KafkaApp) consumeMessages(ctx context.Context) {
    handler := &ConsumerGroupHandler{}
    
    for {
        if err := app.consumerGroup.Consume(ctx, []string{"test-topic"}, handler); err != nil {
            println("Error from consumer group:", err.Error())
        }
        
        if ctx.Err() != nil {
            return
        }
    }
}

func main() {
    gone.NewApp(
        goneKafka.LoaderSyncProducer,
        goneKafka.LoaderAsyncProducer,
        goneKafka.LoadConsumer,
        goneKafka.LoadConsumerGroup,
    ).Run(func(app *KafkaApp) {
        app.Run()
    })
}
```

## Advanced Usage

### Environment Variable Configuration

You can also configure Kafka using environment variables:

```bash
export GONE_KAFKA_DEFAULT='{"addrs":["127.0.0.1:9092"],"groupID":"my-group"}'
```

### Custom Configuration

The component supports all Kafka configuration options provided by the `github.com/IBM/sarama` package, including:

- Producer configuration (retry, acknowledgment, compression, etc.)
- Consumer configuration (offset management, session timeout, etc.)
- Network configuration (timeout, buffer size, etc.)
- Security configuration (SASL, TLS, etc.)
- Management configuration (topic management, partition management, etc.)

For detailed configuration options, please refer to the [Sarama documentation](https://github.com/IBM/sarama).

## Available Loaders

- `LoadConsumer` - Load Kafka consumer
- `LoadConsumerGroup` - Load Kafka consumer group
- `LoaderSyncProducer` - Load synchronous producer
- `LoaderAsyncProducer` - Load asynchronous producer

## Important Notes

1. All client instances will automatically clean up resources when the application shuts down
2. Consumer groups require the `groupID` parameter to be configured
3. Asynchronous producers need to monitor success and error channels
4. It is recommended to configure appropriate retry and timeout parameters in production environments
