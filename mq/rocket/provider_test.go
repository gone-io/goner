package rocket

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"os"
	"reflect"
	"testing"
	"time"

	mq "github.com/apache/rocketmq-clients/golang"
)

func TestConsumerOption_ToOptions(t *testing.T) {

	var Expect = func(awaitDuration1 *time.Duration, subscriptionExpressions1 map[string]*mq.FilterExpression) func() {
		withAwaitDurationExecuted := false
		withSubscriptionExpressionsExecuted := false

		withAwaitDuration = func(awaitDuration time.Duration) mq.SimpleConsumerOption {
			if *awaitDuration1 != awaitDuration {
				t.Errorf("withAwaitDuration() = %v, want %v", awaitDuration1, awaitDuration)
			}
			withAwaitDurationExecuted = true
			return mq.WithAwaitDuration(awaitDuration)
		}
		withSubscriptionExpressions = func(subscriptionExpressions map[string]*mq.FilterExpression) mq.SimpleConsumerOption {
			if !reflect.DeepEqual(subscriptionExpressions1, subscriptionExpressions) {
				t.Errorf("withSubscriptionExpressions() = %v, want %v", subscriptionExpressions1, subscriptionExpressions)
			}
			withSubscriptionExpressionsExecuted = true
			return mq.WithSubscriptionExpressions(subscriptionExpressions)
		}

		return func() {
			if awaitDuration1 != nil && !withAwaitDurationExecuted {
				t.Errorf("withAwaitDuration() not executed")
			}
			if subscriptionExpressions1 != nil && !withSubscriptionExpressionsExecuted {
				t.Errorf("withSubscriptionExpressions() not executed")
			}
		}
	}

	type fields struct {
		AwaitDuration time.Duration
		Expressions   []struct {
			Topic      string `json:"topic"`
			Type       string `json:"type"` // sql or tag
			Expression string `json:"expression"`
		}
	}
	tests := []struct {
		name   string
		setUp  func(fields fields) func()
		fields fields
	}{
		{
			name: "empty option",
			setUp: func(fields fields) func() {
				return Expect(nil, nil)
			},
			fields: fields{
				AwaitDuration: 0,
				Expressions:   nil,
			},
		},
		{
			name: "only await duration",
			setUp: func(fields fields) func() {
				return Expect(&fields.AwaitDuration, nil)
			},
			fields: fields{
				AwaitDuration: 5 * time.Second,
				Expressions:   nil,
			},
		},
		{
			name: "only sql expression",
			setUp: func(fields fields) func() {
				return Expect(nil, map[string]*mq.FilterExpression{
					"test-topic": mq.NewFilterExpressionWithType("age >= 18", mq.SQL92),
				})
			},
			fields: fields{
				AwaitDuration: 0,
				Expressions: []struct {
					Topic      string `json:"topic"`
					Type       string `json:"type"`
					Expression string `json:"expression"`
				}{
					{
						Topic:      "test-topic",
						Type:       "sql",
						Expression: "age >= 18",
					},
				},
			},
		},
		{
			name: "only tag expression",
			setUp: func(fields fields) func() {
				return Expect(nil, map[string]*mq.FilterExpression{
					"test-topic": mq.NewFilterExpressionWithType("TagA || TagB", mq.TAG),
				})
			},
			fields: fields{
				AwaitDuration: 0,
				Expressions: []struct {
					Topic      string `json:"topic"`
					Type       string `json:"type"`
					Expression string `json:"expression"`
				}{
					{
						Topic:      "test-topic",
						Type:       "tag",
						Expression: "TagA || TagB",
					},
				},
			},
		},
		{
			name: "multiple expressions with await duration",
			setUp: func(fields fields) func() {
				return Expect(&fields.AwaitDuration, map[string]*mq.FilterExpression{
					"topic1": mq.NewFilterExpressionWithType("price > 100", mq.SQL92),
					"topic2": mq.NewFilterExpressionWithType("TagX", mq.TAG),
				})
			},
			fields: fields{
				AwaitDuration: 3 * time.Second,
				Expressions: []struct {
					Topic      string `json:"topic"`
					Type       string `json:"type"`
					Expression string `json:"expression"`
				}{
					{
						Topic:      "topic1",
						Type:       "sql",
						Expression: "price > 100",
					},
					{
						Topic:      "topic2",
						Type:       "tag",
						Expression: "TagX",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.setUp(tt.fields)()
			o := &ConsumerOption{
				AwaitDuration: tt.fields.AwaitDuration,
				Expressions:   tt.fields.Expressions,
			}
			_ = o.ToOptions()
		})
	}
}

func TestSendAndReceive(t *testing.T) {
	var topic = "test-topic"
	var maxMessageNum int32 = 16
	var invisibleDuration = time.Second * 20
	var info = "hello `gone`"

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
	"awaitDuration": "10s",
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
						fmt.Println(mv)
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
			time.Sleep(time.Minute)
		})
}
