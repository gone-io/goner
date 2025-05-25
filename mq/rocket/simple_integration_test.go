package rocket

import (
	"context"
	"crypto/rand"
	"fmt"
	mq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/gone-io/gone/v2"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestSendAndReceive(t *testing.T) {
	t.Skipf("skip")

	var topic = "TopicTest"
	var maxMessageNum int32 = 16
	var invisibleDuration = time.Second * 20
	n, _ := rand.Int(rand.Reader, big.NewInt(1000))
	var info = fmt.Sprintf("hello `gone` - %d", n)

	_ = os.Setenv("GONE_ROCKETMQ_DEFAULT", `{
	"Endpoint":"127.0.0.1:8081",
	"ConsumerGroup": "x-test",
	"NameSpace": "xxxx",
	"Credentials": {
		"accessKey":"",
		"accessSecret":"",
		"securityToken":""
	}
}`)

	_ = os.Setenv("GONE_ROCKETMQ_DEFAULT_CONSUMER", fmt.Sprintf(`{
	"awaitDuration": 10000000000,
	"expressions": [{
		"topic": "%s"
	}]
}`, topic))

	_ = os.Setenv("GONE_ROCKETMQ_DEFAULT_PRODUCER", fmt.Sprintf(`{
	"maxAttempts": 3,
	"topics": ["%s"]
}`, topic))

	_ = os.Setenv(mq.CLIENT_LOG_ROOT, "./testdata/logs")
	mq.ResetLogger()

	defer func() {
		_ = os.Unsetenv("GONE_ROCKETMQ_DEFAULT")
		_ = os.Unsetenv("GONE_ROCKETMQ_DEFAULT_CONSUMER")
		_ = os.Unsetenv("GONE_ROCKETMQ_DEFAULT_PRODUCER")
		_ = os.Unsetenv(mq.CLIENT_LOG_ROOT)
	}()

	gone.
		NewApp(LoadConsumer, LoadProducer).
		Run(func(simpleConsumer mq.SimpleConsumer, producer mq.Producer) {
			ch := make(chan struct{})

			go func() {
				for {
					fmt.Println("start receive message")
					mvs, err := simpleConsumer.Receive(context.TODO(), maxMessageNum, invisibleDuration)
					if err != nil {
						fmt.Println(err)
					}

					// ack message
					for _, mv := range mvs {
						_ = simpleConsumer.Ack(context.TODO(), mv)
						fmt.Println("received:", mv)
						if info == string(mv.GetBody()) {
							close(ch)
						}
					}
					fmt.Println("wait a moment")
					fmt.Println()
					time.Sleep(time.Second * 3)
				}
			}()

			go func() {
				msg := &mq.Message{
					Topic: topic,
					Body:  []byte(info),
				}
				// set keys and tag
				msg.SetKeys("a", "b")
				msg.SetTag("ab")
				// send message in sync
				resp, err := producer.Send(context.TODO(), msg)
				if err != nil {
					t.Error(err)
					return
				}
				for i := 0; i < len(resp); i++ {
					fmt.Printf("%#v\n", resp[i])
				}
			}()

			<-ch
		})
}
