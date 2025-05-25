# goner/mq/rabbitmq 组件 和 Gone RabbitMQ 集成

本包为 Gone 应用程序提供 RabbitMQ 集成功能，基于官方 RabbitMQ Go 客户端库提供简单易用的 RabbitMQ 客户端配置和管理。

## 功能特点

- 与 Gone 的依赖注入系统轻松集成
- 支持多个 RabbitMQ 客户端实例
- 提供生产者和消费者接口
- 支持交换机和队列声明
- 支持消息路由和绑定
- 自动资源管理和清理
- 全面的配置选项

## 安装

```bash
gonectl install goner/mq/rabbitmq
gonectl install goner/viper # 可选，用于加载配置文件
```

## 配置

在项目的配置目录中创建 `config/default.yaml` 文件，添加以下 RabbitMQ 配置：

```yaml
rabbitmq:
  default:
    url: "amqp://guest:guest@localhost:5672/"  # RabbitMQ 连接 URL
    # 或者使用单独的连接参数：
    host: "localhost"                          # RabbitMQ 主机
    port: 5672                                 # RabbitMQ 端口
    username: "guest"                          # 认证用户名
    password: "guest"                          # 认证密码
    vhost: "/"                                 # 虚拟主机
    producer:                                  # 生产者配置
      exchange: "my-exchange"                  # 交换机名称
      exchangeType: "direct"                   # 交换机类型 (direct, fanout, topic, headers)
      durable: true                            # 交换机是否持久化
      autoDelete: false                        # 交换机是否自动删除
      internal: false                          # 是否为内部交换机
      noWait: false                            # 是否等待服务器确认
    consumer:                                  # 消费者配置
      queue: "my-queue"                        # 队列名称
      exchange: "my-exchange"                  # 绑定的交换机名称
      routingKey: "my-routing-key"             # 绑定的路由键
      durable: true                            # 队列是否持久化
      autoDelete: false                        # 队列是否自动删除
      exclusive: false                         # 队列是否独占
      noWait: false                            # 是否等待服务器确认
      autoAck: false                           # 是否自动确认消息
      consumer: "my-consumer"                  # 消费者标签
```

### 多客户端配置

你可以配置多个 RabbitMQ 客户端实例：

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

## 使用方法

### 基本用法

#### 生产者

```go
package main

import (
    "github.com/gone-io/gone/v2"
    goneRabbitMQ "github.com/gone-io/goner/mq/rabbitmq"
    amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
    gone.Flag
    producer goneRabbitMQ.IProducer `gone:"*"` // 默认生产者
}

func (p *Producer) SendMessage() {
    // 发送消息
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
        // 处理错误
        return
    }
    
    println("消息发送成功")
}

func main() {
    gone.NewApp(goneRabbitMQ.LoadProducer).Run(func(p *Producer) {
        p.SendMessage()
    })
}
```

#### 消费者

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
    consumer goneRabbitMQ.IConsumer `gone:"*"` // 默认消费者
}

func (c *Consumer) ConsumeMessages(ctx context.Context) {
    // 开始消费消息
    msgs, err := c.consumer.Consume()
    if err != nil {
        log.Fatal("开始消费失败:", err)
    }

    for {
        select {
        case msg := <-msgs:
            // 处理消息
            log.Printf("收到消息: %s", string(msg.Body))

            // 确认消息（如果 autoAck 为 false）
            if err := msg.Ack(false); err != nil {
                log.Printf("确认消息失败: %v", err)
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

### 完整示例

以下是一个完整的生产者和消费者示例，展示了如何在同一个应用程序中使用它们：

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

    // 发送消息
    go app.sendMessages(ctx)

    // 消费消息
    go app.consumeMessages(ctx)

    // 等待中断信号
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
            message := fmt.Sprintf("消息 %d", i)
            err := app.producer.Publish(
                "",     // 路由键
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
                log.Printf("发送消息失败: %v", err)
            } else {
                log.Printf("已发送: %s", message)
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
        log.Fatal("开始消费失败:", err)
    }

    for {
        select {
        case msg := <-msgs:
            log.Printf("收到: %s", string(msg.Body))

            // 确认消息
            if err := msg.Ack(false); err != nil {
                log.Printf("确认消息失败: %v", err)
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

## 环境变量配置

你也可以通过环境变量来配置 RabbitMQ：

```bash
export GONE_RABBITMQ_DEFAULT='{"host":"127.0.0.1","port":5672,"username":"guest","password":"guest"}'
export GONE_RABBITMQ_DEFAULT_PRODUCER='{"exchange":"my-exchange","exchangeType":"direct","autoDelete":true,"durable":false,"noWait":false,"internal":false}'
export GONE_RABBITMQ_DEFAULT_CONSUMER='{"queue":"my-queue","exchange":"my-exchange","routingKey":"my-key","autoDelete":true,"durable":false,"noWait":false,"autoAck":true,"consumer":"my-consumer"}'
```

## 可用的加载器

- `LoadConnection` - 仅加载 RabbitMQ 连接
- `LoadChannel` - 仅加载 RabbitMQ 通道（需要连接）
- `LoadProducer` - 加载 RabbitMQ 生产者（包含连接和通道）
- `LoadConsumer` - 加载 RabbitMQ 消费者（包含连接和通道）
- `LoadAll` - 加载所有 RabbitMQ 组件