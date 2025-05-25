<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/mq/kafka 组件 和 Gone Kafka 集成

本包为 Gone 应用程序提供 Kafka 集成功能，基于 IBM Sarama 库提供简单易用的 Kafka 客户端配置和管理。

## 功能特点

- 与 Gone 的依赖注入系统轻松集成
- 支持多个 Kafka 客户端实例
- 提供同步和异步生产者
- 支持消费者和消费者组
- 自动资源管理和清理
- 全面的配置选项

## 安装

```bash
gonectl install goner/mq/kafka
gonectl install goner/viper # 可选，用于加载配置文件
```

## 配置

在项目的配置目录中创建 `config/default.yaml` 文件，添加以下 Kafka 配置：

```yaml
kafka:
  default:
    addrs: ["127.0.0.1:9092"]     # Kafka 服务器地址列表
    groupID: "my-consumer-group"   # 消费者组 ID
    Producer:
      Return:
        Successes: true            # 是否返回成功消息
        Errors: true               # 是否返回错误消息
      RequiredAcks: 1              # 确认级别
      Retry:
        Max: 3                     # 最大重试次数
    Consumer:
      Offsets:
        AutoCommit:
          Enable: true             # 是否自动提交偏移量
          Interval: 1000           # 自动提交间隔（毫秒）
        Initial: -2                # 初始偏移量（-2: 最早，-1: 最新）
      Group:
        Rebalance:
          Strategy: "range"        # 重平衡策略
```

### 多客户端配置

你可以配置多个 Kafka 客户端实例：

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

## 使用方法

### 基本用法

#### 同步生产者

```go
package main

import (
    "github.com/IBM/sarama"
    "github.com/gone-io/gone/v2"
    goneKafka "github.com/gone-io/goner/mq/kafka"
)

type Producer struct {
    gone.Flag
    syncProducer sarama.SyncProducer `gone:"*"` // 默认同步生产者
}

func (p *Producer) SendMessage() {
    msg := &sarama.ProducerMessage{
        Topic: "my-topic",
        Value: sarama.StringEncoder("Hello, Kafka!"),
    }
    
    partition, offset, err := p.syncProducer.SendMessage(msg)
    if err != nil {
        // 处理错误
        return
    }
    // 消息发送成功
    println("Message sent to partition", partition, "at offset", offset)
}

func main() {
    gone.NewApp(goneKafka.LoaderSyncProducer).Run(func(p *Producer) {
        p.SendMessage()
    })
}
```

#### 异步生产者

```go
type AsyncProducer struct {
    gone.Flag
    asyncProducer sarama.AsyncProducer `gone:"*"` // 默认异步生产者
}

func (p *AsyncProducer) SendMessage() {
    msg := &sarama.ProducerMessage{
        Topic: "my-topic",
        Value: sarama.StringEncoder("Hello, Async Kafka!"),
    }

    // 发送消息到输入通道
    p.asyncProducer.Input() <- msg

    // 监听成功和错误通道
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

#### 消费者

```go
type Consumer struct {
    gone.Flag
    consumer sarama.Consumer `gone:"*"` // 默认消费者
}

func (c *Consumer) ConsumeMessages() {
    partitionConsumer, err := c.consumer.ConsumePartition("my-topic", 0, sarama.OffsetNewest)
    if err != nil {
        // 处理错误
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

#### 消费者组

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
    consumerGroup sarama.ConsumerGroup `gone:"*"` // 默认消费者组
}

func (c *ConsumerGroup) ConsumeMessages(ctx context.Context) {
    handler := &ConsumerGroupHandler{}
    
    for {
        if err := c.consumerGroup.Consume(ctx, []string{"my-topic"}, handler); err != nil {
            // 处理错误
            return
        }
        
        if ctx.Err() != nil {
            return
        }
    }
}
```

### 多客户端示例

```go
type MultiKafkaClient struct {
    gone.Flag
    defaultProducer sarama.SyncProducer `gone:"*"`           // 默认生产者
    cluster1Producer sarama.SyncProducer `gone:"*,cluster1"`  // cluster1 生产者
    cluster2Consumer sarama.Consumer     `gone:"*,cluster2"`  // cluster2 消费者
}

func (m *MultiKafkaClient) UseMultipleClients() {
    // 使用默认集群发送消息
    msg1 := &sarama.ProducerMessage{
        Topic: "topic1",
        Value: sarama.StringEncoder("Message to default cluster"),
    }
    m.defaultProducer.SendMessage(msg1)
    
    // 使用 cluster1 发送消息
    msg2 := &sarama.ProducerMessage{
        Topic: "topic2",
        Value: sarama.StringEncoder("Message to cluster1"),
    }
    m.cluster1Producer.SendMessage(msg2)
    
    // 从 cluster2 消费消息
    partitionConsumer, _ := m.cluster2Consumer.ConsumePartition("topic3", 0, sarama.OffsetNewest)
    defer partitionConsumer.Close()
    
    for message := range partitionConsumer.Messages() {
        println("Received from cluster2:", string(message.Value))
    }
}
```

### 完整示例

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
    
    // 发送同步消息
    go app.sendSyncMessages()
    
    // 发送异步消息
    go app.sendAsyncMessages()
    
    // 消费消息
    go app.consumeMessages(ctx)
    
    // 等待中断信号
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

## 高级用法

### 环境变量配置

你也可以通过环境变量来配置 Kafka：

```bash
export GONE_KAFKA_DEFAULT='{"addrs":["127.0.0.1:9092"],"groupID":"my-group"}'
```

### 自定义配置

组件支持 `github.com/IBM/sarama` 包提供的所有 Kafka 配置选项，包括：

- 生产者配置（重试、确认、压缩等）
- 消费者配置（偏移量管理、会话超时等）
- 网络配置（超时、缓冲区大小等）
- 安全配置（SASL、TLS 等）
- 管理配置（主题管理、分区管理等）

详细的配置选项请参考 [Sarama 文档](https://github.com/IBM/sarama)。

## 可用的加载器

- `LoadConsumer` - 加载 Kafka 消费者
- `LoadConsumerGroup` - 加载 Kafka 消费者组
- `LoaderSyncProducer` - 加载同步生产者
- `LoaderAsyncProducer` - 加载异步生产者

## 注意事项

1. 所有客户端实例都会在应用程序关闭时自动清理资源
2. 消费者组需要配置 `groupID` 参数
3. 异步生产者需要监听成功和错误通道
4. 建议在生产环境中配置适当的重试和超时参数
