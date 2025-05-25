package rabbitmq

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	var (
		exchangeName = "user.register.direct"
		queueName    = "user.register.queue"
		keyName      = "user.register.event"
	)

	t.Run("success", func(t *testing.T) {

		//	_ = os.Setenv("GONE_RABBITMQ_DEFAULT", `{
		//	"host":"127.0.0.1",
		//	"port": 5672,
		//	"username": "guest",
		//	"password": "guest"
		//}`)
		//producer
		_ = os.Setenv("GONE_RABBITMQ_DEFAULT_PRODUCER", `{
	"exchange": "`+exchangeName+`",
	"routingKey": "`+keyName+`",
	"autoDelete": true,
	"durable": false,
	"noWait": false,
	"internal": false
}`)
		_ = os.Setenv("GONE_RABBITMQ_DEFAULT_CONSUMER", `{
	"queue": "`+queueName+`",
	"exchange": "`+exchangeName+`",
	"routingKey": "`+keyName+`",
	"autoDelete": true,
	"durable": false,
	"noWait": false,
	"autoAck": true,
	"consumer": "my-test"
}`)
		defer func() {
			//_ = os.Unsetenv("GONE_RABBITMQ_DEFAULT")
			_ = os.Unsetenv("GONE_RABBITMQ_DEFAULT_PRODUCER")
			_ = os.Unsetenv("GONE_RABBITMQ_DEFAULT_CONSUMER")
		}()

		gone.
			NewApp(LoadAll).
			Run(func(p IProducer, c IConsumer) {
				assert.NotNil(t, p)
				assert.NotNil(t, c)

				info := "hello gone"
				ch := make(chan struct{})

				go func() {
					consume, err := c.Consume()
					assert.NoError(t, err)
					for {
						select {
						case msg := <-consume:
							if info == string(msg.Body) {
								close(ch)
							}
						}
					}
				}()

				err := p.Publish(keyName, false, false, amqp.Publishing{
					Headers:      nil,
					ContentType:  "text/plain",
					DeliveryMode: amqp.Persistent,
					Timestamp:    time.Now(),
					Body:         []byte(info),
				})
				assert.NoError(t, err)
				<-ch
			})
	})

	t.Run("Exchange is not set", func(t *testing.T) {
		_ = os.Setenv("GONE_RABBITMQ_DEFAULT", `{
			"url":"`+fmt.Sprintf("amqp://%s:%s@%s:%d%s", "guest", "guest", "127.0.0.1", 5672, "/")+`",
			"host":"127.0.0.1",
			"port": 5672,
			"username": "guest",
			"password": "guest"
		}`)
		defer func() {
			_ = os.Unsetenv("GONE_RABBITMQ_DEFAULT")
		}()

		err := gone.SafeExecute(func() error {
			gone.NewApp(LoadAll).
				Run(func(p IProducer) {
				})
			return nil
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exchange name is required")
	})

	t.Run("queue name is required", func(t *testing.T) {
		err := gone.SafeExecute(func() error {
			gone.NewApp(LoadAll).
				Run(func(p IConsumer) {
				})
			return nil
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "queue name is required")
	})
	t.Run("exchange name is required", func(t *testing.T) {
		_ = os.Setenv("GONE_RABBITMQ_DEFAULT_CONSUMER", `{
	"queue": "`+queueName+`",
	"routingKey": "`+keyName+`",
	"autoDelete": true,
	"durable": false,
	"noWait": false,
	"autoAck": true,
	"consumer": "my-test"
}`)

		err := gone.SafeExecute(func() error {
			gone.NewApp(LoadAll).
				Run(func(p IConsumer) {
				})
			return nil
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exchange name is required")
	})
}
