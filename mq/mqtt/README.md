<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# goner/mq/mqtt Component and Gone MQTT Integration

This package provides MQTT integration functionality for Gone applications, offering simple and easy-to-use MQTT client configuration and management based on the Eclipse Paho MQTT library.

## Features

- Easy integration with Gone's dependency injection system
- Support for multiple MQTT client instances
- Automatic connection management and cleanup
- Comprehensive configuration options
- Support for publish/subscribe messaging patterns
- Built-in reconnection handling

## Installation

```bash
gonectl install goner/mq/mqtt
gonectl install goner/viper # Optional, for loading configuration files
```

## Configuration

Create a `config/default.yaml` file in your project's configuration directory and add the following MQTT configuration:

```yaml
mqtt:
  default:
    brokers: ["tcp://localhost:1883"]  # List of MQTT broker addresses
    ClientID: "my-client-id"           # Client ID for MQTT connection
    Username: "username"               # Username for authentication (optional)
    Password: "password"               # Password for authentication (optional)
    CleanSession: true                 # Clean session flag
    KeepAlive: 60                      # Keep alive interval in seconds
    ConnectTimeout: 30                 # Connection timeout in seconds
    AutoReconnect: true                # Enable automatic reconnection
    MaxReconnectInterval: 10           # Maximum reconnection interval in seconds
    MessageChannelDepth: 100           # Message channel depth
```

### Multiple Client Configuration

You can configure multiple MQTT client instances:

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

## Usage

### Basic Usage

#### Publisher

```go
package main

import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/gone-io/gone/v2"
    goneMqtt "github.com/gone-io/goner/mq/mqtt"
)

type Publisher struct {
    gone.Flag
    client mqtt.Client `gone:"*"` // Default MQTT client
}

func (p *Publisher) PublishMessage() {
    topic := "sensors/temperature"
    payload := "25.5"
    
    token := p.client.Publish(topic, 1, false, payload)
    token.Wait()
    if err := token.Error(); err != nil {
        // Handle error
        return
    }
    // Message published successfully
    println("Message published to", topic)
}

func main() {
    gone.NewApp(goneMqtt.Load).Run(func(p *Publisher) {
        p.PublishMessage()
    })
}
```

#### Subscriber

```go
type Subscriber struct {
    gone.Flag
    client mqtt.Client `gone:"*"` // Default MQTT client
}

func (s *Subscriber) SubscribeToTopic() {
    topic := "sensors/+"  // Subscribe to all sensor topics
    
    token := s.client.Subscribe(topic, 1, s.messageHandler)
    token.Wait()
    if err := token.Error(); err != nil {
        // Handle error
        return
    }
    println("Subscribed to", topic)
}

func (s *Subscriber) messageHandler(client mqtt.Client, msg mqtt.Message) {
    println("Received message:")
    println("Topic:", msg.Topic())
    println("Payload:", string(msg.Payload()))
    println("QoS:", msg.Qos())
    println("Retained:", msg.Retained())
}
```

#### Pub/Sub Example

```go
type PubSubApp struct {
    gone.Flag
    client mqtt.Client `gone:"*"`
}

func (app *PubSubApp) Run() {
    // Subscribe to topics
    app.subscribeToTopics()
    
    // Publish messages
    app.publishMessages()
    
    // Keep the application running
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
        println("Subscribe error:", err.Error())
        return
    }
    println("Subscribed to multiple topics")
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
            println("Publish error:", err.Error())
        } else {
            println("Published:", msg.topic, "=", msg.payload)
        }
    }
}

func (app *PubSubApp) messageHandler(client mqtt.Client, msg mqtt.Message) {
    println("[RECEIVED]", msg.Topic(), ":", string(msg.Payload()))
}
```

### Multiple Clients Example

```go
type MultiMqttClient struct {
    gone.Flag
    defaultClient mqtt.Client `gone:"*"`           // Default client
    sensorClient  mqtt.Client `gone:"*,sensor"`    // Sensor client
    actuatorClient mqtt.Client `gone:"*,actuator"`  // Actuator client
}

func (m *MultiMqttClient) UseMultipleClients() {
    // Publish sensor data using sensor client
    token := m.sensorClient.Publish("sensors/temperature", 1, false, "23.5")
    token.Wait()
    
    // Subscribe to commands using actuator client
    token = m.actuatorClient.Subscribe("commands/+", 1, func(client mqtt.Client, msg mqtt.Message) {
        println("Command received:", string(msg.Payload()))
    })
    token.Wait()
    
    // Use default client for general messaging
    token = m.defaultClient.Publish("status/online", 0, true, "true")
    token.Wait()
}
```

### Complete Example

```go
package main

import (
    "context"
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
    
    // Subscribe to topics
    go app.subscribe()
    
    // Publish messages periodically
    go app.publishPeriodically(ctx)
    
    // Wait for interrupt signal
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    <-signals
    
    println("Shutting down...")
}

func (app *MqttApp) subscribe() {
    token := app.client.Subscribe("test/+", 1, func(client mqtt.Client, msg mqtt.Message) {
        println("[RECEIVED]", msg.Topic(), ":", string(msg.Payload()))
    })
    token.Wait()
    if err := token.Error(); err != nil {
        println("Subscribe error:", err.Error())
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
            payload := fmt.Sprintf("Message %d", counter)
            token := app.client.Publish("test/messages", 1, false, payload)
            token.Wait()
            if err := token.Error(); err != nil {
                println("Publish error:", err.Error())
            } else {
                println("[PUBLISHED]", payload)
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

## Advanced Usage

### Environment Variable Configuration

You can also configure MQTT using environment variables:

```bash
export GONE_MQTT_DEFAULT='{"brokers":["tcp://localhost:1883"],"ClientID":"my-client"}'
```

### Custom Configuration

The component supports all MQTT configuration options provided by the `github.com/eclipse/paho.mqtt.golang` package, including:

- Connection configuration (brokers, client ID, credentials, etc.)
- Session configuration (clean session, keep alive, timeouts, etc.)
- Reconnection configuration (auto reconnect, retry intervals, etc.)
- Message configuration (QoS levels, retained messages, etc.)
- TLS/SSL configuration for secure connections
- Will message configuration for last will and testament

For detailed configuration options, please refer to the [Paho MQTT documentation](https://github.com/eclipse/paho.mqtt.golang).

## Available Loader

- `Load` - Load MQTT client

## Important Notes

1. All client instances will automatically disconnect when the application shuts down
2. Client IDs should be unique across all clients connecting to the same broker
3. QoS levels determine message delivery guarantees (0: at most once, 1: at least once, 2: exactly once)
4. Retained messages are stored by the broker and delivered to new subscribers
5. It is recommended to handle connection lost scenarios in production environments
6. Use appropriate QoS levels based on your application requirements