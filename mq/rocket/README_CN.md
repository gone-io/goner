<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/mq/rocket 组件 和 Gone RocketMQ 集成

本包为 Gone 应用程序提供 RocketMQ 集成功能，基于 Apache RocketMQ Go 客户端库提供简单易用的 RocketMQ 客户端配置和管理。

## 功能特点

- 与 Gone 的依赖注入系统轻松集成
- 支持多个 RocketMQ 客户端实例
- 提供生产者和简单消费者
- 支持事务消息
- 支持 SQL 和 Tag 过滤表达式
- 自动资源管理和清理
- 全面的配置选项

## 安装

```bash
gonectl install goner/mq/rocket
gonectl install goner/viper # 可选，用于加载配置文件
```

## 配置

在项目的配置目录中创建 `config/default.yaml` 文件，添加以下 RocketMQ 配置：

```yaml
rocketmq:
  default:
    Endpoint: "127.0.0.1:8081"          # RocketMQ 代理服务器地址
    ConsumerGroup: "my-consumer-group"    # 消费者组 ID
    NameSpace: ""                        # 命名空间（可选）
    Credentials:                          # 认证信息（可选）
      accessKey: ""
      accessSecret: ""
      securityToken: ""
    consumer:                             # 消费者配置
      awaitDuration: 10000000000          # 等待时间（纳秒）
      expressions:                        # 订阅表达式
        - topic: "my-topic"
          type: "tag"                      # 过滤类型：tag 或 sql
          expression: "*"                  # 过滤表达式
    producer:                             # 生产者配置
      maxAttempts: 3                      # 最大重试次数
      topics: ["my-topic"]               # 预设主题列表
```

### 多客户端配置

你可以配置多个 RocketMQ 客户端实例：

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

## 使用方法

### 基本用法

#### 生产者

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
    producer mq.Producer `gone:"*"` // 默认生产者
}

func (p *Producer) SendMessage() {
    msg := &mq.Message{
        Topic: "my-topic",
        Body:  []byte("Hello, RocketMQ!"),
    }
    
    // 设置消息属性
    msg.SetKeys("key1", "key2")
    msg.SetTag("important")
    
    // 发送消息
    resp, err := p.producer.Send(context.TODO(), msg)
    if err != nil {
        // 处理错误
        return
    }
    
    // 处理发送结果
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

#### 简单消费者

```go
type Consumer struct {
    gone.Flag
    consumer mq.SimpleConsumer `gone:"*"` // 默认简单消费者
}

func (c *Consumer) ConsumeMessages() {
    ctx := context.Background()
    maxMessageNum := int32(16)
    invisibleDuration := time.Second * 20
    
    for {
        // 接收消息
        messages, err := c.consumer.Receive(ctx, maxMessageNum, invisibleDuration)
        if err != nil {
            println("Failed to receive messages:", err.Error())
            continue
        }
        
        // 处理消息
        for _, message := range messages {
            println("Received message:", string(message.GetBody()))
            
            // 确认消息
            if err := c.consumer.Ack(ctx, message); err != nil {
                println("Failed to ack message:", err.Error())
            }
        }
        
        time.Sleep(time.Second)
    }
}
```

### 多客户端示例

```go
type MultiRocketMQClient struct {
    gone.Flag
    defaultProducer mq.Producer       `gone:"*"`           // 默认生产者
    cluster1Producer mq.Producer      `gone:"*,cluster1"`  // cluster1 生产者
    cluster2Consumer mq.SimpleConsumer `gone:"*,cluster2"`  // cluster2 消费者
}

func (m *MultiRocketMQClient) UseMultipleClients() {
    ctx := context.Background()
    
    // 使用默认集群发送消息
    msg1 := &mq.Message{
        Topic: "topic1",
        Body:  []byte("Message to default cluster"),
    }
    m.defaultProducer.Send(ctx, msg1)
    
    // 使用 cluster1 发送消息
    msg2 := &mq.Message{
        Topic: "topic2",
        Body:  []byte("Message to cluster1"),
    }
    m.cluster1Producer.Send(ctx, msg2)
    
    // 从 cluster2 消费消息
    messages, _ := m.cluster2Consumer.Receive(ctx, 10, time.Second*20)
    for _, message := range messages {
        println("Received from cluster2:", string(message.GetBody()))
        m.cluster2Consumer.Ack(ctx, message)
    }
}
```

### 过滤表达式示例

#### Tag 过滤

```yaml
rocketmq:
  default:
    consumer:
      expressions:
        - topic: "order-topic"
          type: "tag"
          expression: "important || urgent"  # 接收 important 或 urgent 标签的消息
```

#### SQL 过滤

```yaml
rocketmq:
  default:
    consumer:
      expressions:
        - topic: "user-topic"
          type: "sql"
          expression: "age >= 18 AND region = 'US'"  # SQL 表达式过滤
```

### 完整示例

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
    
    // 发送消息
    go app.sendMessages(ctx)
    
    // 消费消息
    go app.consumeMessages(ctx)
    
    // 等待中断信号
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
        
        // 设置消息属性
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
                
                // 确认消息
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

## 高级用法

### 环境变量配置

你也可以通过环境变量来配置 RocketMQ：

```bash
export GONE_ROCKETMQ_DEFAULT='{"Endpoint":"127.0.0.1:8081","ConsumerGroup":"my-group"}'
export GONE_ROCKETMQ_DEFAULT_CONSUMER='{"awaitDuration":10000000000,"expressions":[{"topic":"my-topic"}]}'
export GONE_ROCKETMQ_DEFAULT_PRODUCER='{"maxAttempts":3,"topics":["my-topic"]}'
```

### 事务消息支持

组件支持事务消息，你可以通过注册 `TransactionChecker` 来处理事务回查：

```go
type MyTransactionChecker struct {
    gone.Flag
}

func (c *MyTransactionChecker) Check(msg *mq.MessageView) mq.TransactionResolution {
    // 实现事务状态检查逻辑
    return mq.COMMIT
}

// 注册事务检查器
func (c *MyTransactionChecker) GetGoneId() string {
    return "rocketmq.default" // 对应配置中的客户端名称
}
```

### 自定义配置

组件支持 `github.com/apache/rocketmq-clients/golang/v5` 包提供的所有 RocketMQ 配置选项，包括：

- 生产者配置（重试次数、主题列表等）
- 消费者配置（等待时间、过滤表达式等）
- 网络配置（端点、超时等）
- 安全配置（认证信息等）
- 命名空间配置

详细的配置选项请参考 [RocketMQ Go 客户端文档](https://github.com/apache/rocketmq-clients/tree/master/golang)。

## 可用的加载器

- `LoadConsumer` - 加载 RocketMQ 简单消费者
- `LoadProducer` - 加载 RocketMQ 生产者

## 注意事项

1. 所有客户端实例都会在应用程序关闭时自动清理资源
2. 消费者需要配置 `ConsumerGroup` 参数
3. 建议在生产环境中配置适当的重试和超时参数
4. 使用 SQL 过滤表达式时，需要确保 RocketMQ 服务器支持该功能
5. 事务消息需要实现 `TransactionChecker` 接口来处理事务回查
6. 消费者的 `awaitDuration` 参数以纳秒为单位
