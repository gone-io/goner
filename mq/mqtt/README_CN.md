<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# goner/mq/mqtt 组件 和 Gone MQTT 集成

本包为 Gone 应用程序提供 MQTT 集成功能，基于 Eclipse Paho MQTT 库提供简单易用的 MQTT 客户端配置和管理。

## 功能特点

- 与 Gone 的依赖注入系统轻松集成
- 支持多个 MQTT 客户端实例
- 自动连接管理和清理
- 全面的配置选项
- 支持发布/订阅消息模式
- 内置重连处理机制

## 安装

```bash
gonectl install goner/mq/mqtt
gonectl install goner/viper # 可选，用于加载配置文件
```

## 配置

在项目的配置目录中创建 `config/default.yaml` 文件，添加以下 MQTT 配置：

```yaml
mqtt:
  default:
    brokers: ["tcp://localhost:1883"]  # MQTT 代理服务器地址列表
    ClientID: "my-client-id"           # MQTT 连接的客户端 ID
    Username: "username"               # 认证用户名（可选）
    Password: "password"               # 认证密码（可选）
    CleanSession: true                 # 清理会话标志
    KeepAlive: 60                      # 保持连接间隔（秒）
    ConnectTimeout: 30                 # 连接超时时间（秒）
    AutoReconnect: true                # 启用自动重连
    MaxReconnectInterval: 10           # 最大重连间隔（秒）
    MessageChannelDepth: 100           # 消息通道深度
```

### 多客户端配置

你可以配置多个 MQTT 客户端实例：

```yaml
mqtt:
  default:
    brokers: ["tcp://localhost:1883"]
    ClientID: "default-client"
  sensor:
    brokers: ["tcp://sensor.example.com:1883"]
    ClientID: "sensor-client"
    Username: "sensor_user"
    Password: "sensor_pass"
  actuator:
    brokers: ["tcp://actuator.example.com:1883"]
    ClientID: "actuator-client"
```

## 使用方法

### 基本用法

#### 发布者

```go
package main

import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/gone-io/gone/v2"
    goneMqtt "github.com/gone-io/goner/mq/mqtt"
)

type Publisher struct {
    gone.Flag
    client mqtt.Client `gone:"*"` // 默认 MQTT 客户端
}

func (p *Publisher) PublishMessage() {
    topic := "sensors/temperature"
    payload := "25.5"
    
    token := p.client.Publish(topic, 1, false, payload)
    token.Wait()
    if err := token.Error(); err != nil {
        // 处理错误
        return
    }
    // 消息发布成功
    println("消息已发布到", topic)
}

func main() {
    gone.NewApp(goneMqtt.Load).Run(func(p *Publisher) {
        p.PublishMessage()
    })
}
```

#### 订阅者

```go
type Subscriber struct {
    gone.Flag
    client mqtt.Client `gone:"*"` // 默认 MQTT 客户端
}

func (s *Subscriber) SubscribeToTopic() {
    topic := "sensors/+"  // 订阅所有传感器主题
    
    token := s.client.Subscribe(topic, 1, s.messageHandler)
    token.Wait()
    if err := token.Error(); err != nil {
        // 处理错误
        return
    }
    println("已订阅", topic)
}

func (s *Subscriber) messageHandler(client mqtt.Client, msg mqtt.Message) {
    println("收到消息:")
    println("主题:", msg.Topic())
    println("负载:", string(msg.Payload()))
    println("QoS:", msg.Qos())
    println("保留:", msg.Retained())
}
```

#### 发布/订阅示例

```go
type PubSubApp struct {
    gone.Flag
    client mqtt.Client `gone:"*"`
}

func (app *PubSubApp) Run() {
    // 订阅主题
    app.subscribeToTopics()
    
    // 发布消息
    app.publishMessages()
    
    // 保持应用程序运行
    select {}
}

func (app *PubSubApp) subscribeToTopics() {
    topics := map[string]byte{
        "home/livingroom/temperature": 1,
        "home/bedroom/humidity":       1,
        "home/kitchen/motion":         0,
    }
    
    token := app.client.SubscribeMultiple(topics, app.messageHandler)
    token.Wait()
    if err := token.Error(); err != nil {
        println("订阅错误:", err.Error())
        return
    }
    println("已订阅多个主题")
}

func (app *PubSubApp) publishMessages() {
    messages := []struct {
        topic   string
        payload string
        qos     byte
    }{
        {"home/livingroom/temperature", "22.5", 1},
        {"home/bedroom/humidity", "65", 1},
        {"home/kitchen/motion", "detected", 0},
    }
    
    for _, msg := range messages {
        token := app.client.Publish(msg.topic, msg.qos, false, msg.payload)
        token.Wait()
        if err := token.Error(); err != nil {
            println("发布错误:", err.Error())
        } else {
            println("已发布:", msg.topic, "=", msg.payload)
        }
    }
}

func (app *PubSubApp) messageHandler(client mqtt.Client, msg mqtt.Message) {
    println("[收到消息]", msg.Topic(), ":", string(msg.Payload()))
}
```

### 多客户端示例

```go
type MultiMqttClient struct {
    gone.Flag
    defaultClient  mqtt.Client `gone:"*"`           // 默认客户端
    sensorClient   mqtt.Client `gone:"*,sensor"`    // 传感器客户端
    actuatorClient mqtt.Client `gone:"*,actuator"`  // 执行器客户端
}

func (m *MultiMqttClient) UseMultipleClients() {
    // 使用传感器客户端发布传感器数据
    token := m.sensorClient.Publish("sensors/temperature", 1, false, "23.5")
    token.Wait()
    
    // 使用执行器客户端订阅命令
    token = m.actuatorClient.Subscribe("commands/+", 1, func(client mqtt.Client, msg mqtt.Message) {
        println("收到命令:", string(msg.Payload()))
    })
    token.Wait()
    
    // 使用默认客户端进行一般消息传递
    token = m.defaultClient.Publish("status/online", 0, true, "true")
    token.Wait()
}
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
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/gone-io/gone/v2"
    goneMqtt "github.com/gone-io/goner/mq/mqtt"
)

type MqttApp struct {
    gone.Flag
    client mqtt.Client `gone:"*"`
}

func (app *MqttApp) Run() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 订阅主题
    go app.subscribe()
    
    // 定期发布消息
    go app.publishPeriodically(ctx)
    
    // 等待中断信号
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    <-signals
    
    println("正在关闭...")
}

func (app *MqttApp) subscribe() {
    token := app.client.Subscribe("test/+", 1, func(client mqtt.Client, msg mqtt.Message) {
        println("[收到消息]", msg.Topic(), ":", string(msg.Payload()))
    })
    token.Wait()
    if err := token.Error(); err != nil {
        println("订阅错误:", err.Error())
    }
}

func (app *MqttApp) publishPeriodically(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    counter := 0
    for {
        select {
        case <-ticker.C:
            counter++
            payload := fmt.Sprintf("消息 %d", counter)
            token := app.client.Publish("test/messages", 1, false, payload)
            token.Wait()
            if err := token.Error(); err != nil {
                println("发布错误:", err.Error())
            } else {
                println("[已发布]", payload)
            }
        case <-ctx.Done():
            return
        }
    }
}

func main() {
    gone.NewApp(goneMqtt.Load).Run(func(app *MqttApp) {
        app.Run()
    })
}
```

## 高级用法

### 环境变量配置

你也可以通过环境变量来配置 MQTT：

```bash
export GONE_MQTT_DEFAULT='{"brokers":["tcp://localhost:1883"],"ClientID":"my-client"}'
```

### 自定义配置

组件支持 `github.com/eclipse/paho.mqtt.golang` 包提供的所有 MQTT 配置选项，包括：

- 连接配置（代理服务器、客户端 ID、凭据等）
- 会话配置（清理会话、保持连接、超时等）
- 重连配置（自动重连、重试间隔等）
- 消息配置（QoS 级别、保留消息等）
- TLS/SSL 配置用于安全连接
- 遗嘱消息配置

详细的配置选项请参考 [Paho MQTT 文档](https://github.com/eclipse/paho.mqtt.golang)。

## 可用的加载器

- `Load` - 加载 MQTT 客户端

## 注意事项

1. 所有客户端实例都会在应用程序关闭时自动断开连接
2. 客户端 ID 在连接到同一代理的所有客户端中应该是唯一的
3. QoS 级别决定消息传递保证（0：最多一次，1：至少一次，2：恰好一次）
4. 保留消息由代理存储并传递给新订阅者
5. 建议在生产环境中处理连接丢失场景
6. 根据应用程序要求使用适当的 QoS 级别
