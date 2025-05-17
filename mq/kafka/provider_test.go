package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"os/signal"
	"sync"
	"testing"
)

func TestConf_ReadFromConfigure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	configure := gone.NewMockConfigure(controller)
	_ = configure.EXPECT().
		Get("kafka.default", gomock.Any(), "").
		Do(func(key string, val interface{}, defaultVal string) {
			conf := val.(*Conf)
			conf.Addrs = []string{"127.0.0.1:9092"}
			conf.Config = &sarama.Config{
				ChannelBufferSize: 512,
			}
		}).
		Return(nil)

	var conf Conf
	address, config := conf.ReadFromConfigure("", configure)
	assert.Equal(t, "127.0.0.1:9092", address[0])

	assert.Equal(t, 512, config.ChannelBufferSize)
	assert.Equal(t, 5, config.Admin.Retry.Max)
}

type consumerHandler struct {
	wants  []string
	locker sync.Mutex
	ch     chan string
}

func (h *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.process(msg)
		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}

func (h *consumerHandler) process(msg *sarama.ConsumerMessage) {
	h.locker.Lock()
	defer h.locker.Unlock()
	fmt.Printf("Received message: key=%s, value=%s, partition=%d, offset=%d\n", string(msg.Key), string(msg.Value), msg.Partition, msg.Offset)
	if len(h.wants) > 0 {
		r := -1
		for i, want := range h.wants {
			if want == string(msg.Value) {
				r = i
				break
			}
		}
		if r > -1 {
			//删除 r
			h.wants = append(h.wants[:r], h.wants[r+1:]...)
		}
		if len(h.wants) == 0 {
			close(h.ch)
		}
	}
}

func TestSendAndReceive(t *testing.T) {
	conf := `{
	"groupID": "default",
	"addrs": ["localhost:9092"],
	"Producer": {
		"Return": {
			"Successes": true
		}
	},
	"Consumer": {
		"Offsets": {
			"AutoCommit": {
				"Enable": true
			},
			"Initial": -2
		}
	}
}`

	_ = os.Setenv("GONE_KAFKA_DEFAULT", conf)
	defer func() {
		_ = os.Unsetenv("GONE_KAFKA_DEFAULT")
	}()

	var topic = "my-topic"
	var info1 = "hello"
	var info2 = "gone"

	gone.
		NewApp(LoadConsumer, LoaderSyncProducer, LoaderAsyncProducer, LoadConsumerGroup).
		Run(func(syncProducer sarama.SyncProducer, aSyncProducer sarama.AsyncProducer, consumer sarama.Consumer, client sarama.ConsumerGroup) {
			ctx, cancel := context.WithCancel(context.Background())
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)

			ch := make(chan string)

			go func() {
				for {
					err := client.Consume(ctx, []string{topic}, &consumerHandler{
						wants: []string{info1, info2},
						ch:    ch,
					})
					if err != nil {
						t.Errorf("Error from consumer: %s", err)
					}
					select {
					case <-signals:
						cancel()
						return
					case <-ch:
						cancel()
						return
					}
				}
			}()

			go func() {
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(info1),
				}

				partition, offset, err := syncProducer.SendMessage(msg)
				assert.Nil(t, err)
				fmt.Printf("send:%s,%#v, %#v\n", info1, partition, offset)

				_, err = consumer.Topics()
				assert.Nil(t, err)
			}()

			go func() {
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(info2),
				}
				aSyncProducer.Input() <- msg

				fmt.Printf("send:%s\n", info2)
			}()

			select {
			case <-signals:
				cancel()
				return
			case <-ch:
				cancel()
				return
			}
		})
}
