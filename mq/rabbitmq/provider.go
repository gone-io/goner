package rabbitmq

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Config represents RabbitMQ connection configuration
type Config struct {
	URL      string `json:"url"`      // RabbitMQ connection URL
	Host     string `json:"host"`     // RabbitMQ host
	Port     int    `json:"port"`     // RabbitMQ port
	Username string `json:"username"` // Username for authentication
	Password string `json:"password"` // Password for authentication
	Vhost    string `json:"vhost"`    // Virtual host
}

// GetConnectionURL returns the connection URL for RabbitMQ
func (c *Config) GetConnectionURL() string {
	if c.URL != "" {
		return c.URL
	}

	host := c.Host
	if host == "" {
		host = "localhost"
	}

	port := c.Port
	if port == 0 {
		port = 5672
	}

	username := c.Username
	if username == "" {
		username = "guest"
	}

	password := c.Password
	if password == "" {
		password = "guest"
	}

	vhost := c.Vhost
	if vhost == "" {
		vhost = "/"
	}

	return fmt.Sprintf("amqp://%s:%s@%s:%d%s", username, password, host, port, vhost)
}

// ProducerOption represents producer configuration options
type ProducerOption struct {
	Exchange     string `json:"exchange"`     // Exchange name
	ExchangeType string `json:"exchangeType"` // Exchange type (direct, fanout, topic, headers)
	Durable      bool   `json:"durable"`      // Whether exchange is durable
	AutoDelete   bool   `json:"autoDelete"`   // Whether exchange auto-deletes
	Internal     bool   `json:"internal"`     // Whether exchange is internal
	NoWait       bool   `json:"noWait"`       // Whether to wait for server confirmation
}

// ConsumerOption represents consumer configuration options
type ConsumerOption struct {
	Queue      string `json:"queue"`      // Queue name
	Exchange   string `json:"exchange"`   // Exchange name to bind to
	RoutingKey string `json:"routingKey"` // Routing key for binding
	Durable    bool   `json:"durable"`    // Whether queue is durable
	AutoDelete bool   `json:"autoDelete"` // Whether queue auto-deletes
	Exclusive  bool   `json:"exclusive"`  // Whether queue is exclusive
	NoWait     bool   `json:"noWait"`     // Whether to wait for server confirmation
	AutoAck    bool   `json:"autoAck"`    // Whether to auto-acknowledge messages
	Consumer   string `json:"consumer"`   // Consumer tag
}

var dialFunc = amqp.Dial

// ProvideConnection provides a RabbitMQ connection
func ProvideConnection(tagConf string, param struct {
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (*amqp.Connection, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("rabbitmq.%s", name)

	var conf Config
	err := param.configure.Get(confKey, &conf, "")
	g.PanicIfErr(gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", confKey)))

	conn, err := dialFunc(conf.GetConnectionURL())
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not create rabbitmq connection"))

	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, conn.Close(), "close rabbitmq connection failed")
	})

	return conn, nil
}

// ProvideChannel provides a RabbitMQ channel
func ProvideChannel(tagConf string, param struct {
	connection *amqp.Connection `gone:"*"`
	beforeStop gone.BeforeStop  `gone:"*"`
	logger     gone.Logger      `gone:"*"`
}) (*amqp.Channel, error) {
	ch, err := param.connection.Channel()
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not create rabbitmq channel"))

	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, ch.Close(), "close rabbitmq channel failed")
	})
	return ch, nil
}

type IProducer interface {
	Publish(key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// Producer wraps RabbitMQ channel for publishing messages
type Producer struct {
	channel *amqp.Channel
	option  ProducerOption
}

// NewProducer creates a new Producer instance
func NewProducer(channel *amqp.Channel, option ProducerOption) *Producer {
	return &Producer{
		channel: channel,
		option:  option,
	}
}

// DeclareExchange declares an exchange
func (p *Producer) DeclareExchange() error {
	if p.option.Exchange == "" {
		return gone.NewInnerError("exchange name is required", gone.InjectError)
	}

	exchangeType := p.option.ExchangeType
	if exchangeType == "" {
		exchangeType = "direct"
	}

	return p.channel.ExchangeDeclare(
		p.option.Exchange,
		exchangeType,
		p.option.Durable,
		p.option.AutoDelete,
		p.option.Internal,
		p.option.NoWait,
		nil,
	)
}

// Publish publishes a message to the exchange
func (p *Producer) Publish(routingKey string, mandatory, immediate bool, msg amqp.Publishing) error {
	return p.channel.Publish(
		p.option.Exchange,
		routingKey,
		mandatory, // mandatory
		immediate, // immediate
		msg,
	)
}

type IConsumer interface {
	Consume() (<-chan amqp.Delivery, error)
}

// Consumer wraps RabbitMQ channel for consuming messages
type Consumer struct {
	channel *amqp.Channel
	option  ConsumerOption
}

// NewConsumer creates a new Consumer instance
func NewConsumer(channel *amqp.Channel, option ConsumerOption) *Consumer {
	return &Consumer{
		channel: channel,
		option:  option,
	}
}

// DeclareQueue declares a queue
func (c *Consumer) DeclareQueue() error {
	if c.option.Queue == "" {
		return gone.NewInnerError("queue name is required", gone.InjectError)
	}

	_, err := c.channel.QueueDeclare(
		c.option.Queue,
		c.option.Durable,
		c.option.AutoDelete,
		c.option.Exclusive,
		c.option.NoWait,
		nil,
	)
	return err
}

// BindQueue binds the queue to an exchange
func (c *Consumer) BindQueue() error {
	if c.option.Exchange == "" {
		return gone.NewInnerError("exchange name is required", gone.InjectError)
	}

	return c.channel.QueueBind(
		c.option.Queue,
		c.option.RoutingKey,
		c.option.Exchange,
		c.option.NoWait,
		nil,
	)
}

// Consume starts consuming messages from the queue
func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		c.option.Queue,
		c.option.Consumer,
		c.option.AutoAck,
		c.option.Exclusive,
		false, // no-local
		c.option.NoWait,
		nil,
	)
}

// ProvideProducer provides a RabbitMQ producer
func ProvideProducer(tagConf string, param struct {
	channel    *amqp.Channel   `gone:"*"`
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (*Producer, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	if name == "" {
		name = "default"
	}

	var option ProducerOption
	_ = param.configure.Get(fmt.Sprintf("rabbitmq.%s.producer", name), &option, "")

	producer := NewProducer(param.channel, option)

	// Declare exchange if configured
	err := producer.DeclareExchange()
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not declare exchange"))

	return producer, nil
}

// ProvideConsumer provides a RabbitMQ consumer
func ProvideConsumer(tagConf string, param struct {
	channel    *amqp.Channel   `gone:"*"`
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (*Consumer, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	if name == "" {
		name = "default"
	}

	var option ConsumerOption
	_ = param.configure.Get(fmt.Sprintf("rabbitmq.%s.consumer", name), &option, "")

	consumer := NewConsumer(param.channel, option)

	// Declare queue if configured
	err := consumer.DeclareQueue()
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not declare queue"))

	// Bind queue to exchange if configured
	err = consumer.BindQueue()
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not bind queue"))

	return consumer, nil
}
