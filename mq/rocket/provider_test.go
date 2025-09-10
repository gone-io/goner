package rocket

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
	"time"

	mq "github.com/apache/rocketmq-clients/golang/v5"
)

func TestProducerOption_ToOptions(t *testing.T) {
	var Expect = func(maxAttempts1 *int32, topics1 []string) func() {
		withMaxAttemptsExecuted := false
		withTopicsExecuted := false

		withMaxAttempts = func(maxAttempts int32) mq.ProducerOption {
			if *maxAttempts1 != maxAttempts {
				t.Errorf("withMaxAttempts() = %v, want %v", maxAttempts1, maxAttempts)
			}
			withMaxAttemptsExecuted = true
			return mq.WithMaxAttempts(maxAttempts)
		}
		withTopics = func(topics ...string) mq.ProducerOption {
			if !reflect.DeepEqual(topics1, topics) {
				t.Errorf("withTopics() = %v, want %v", topics1, topics)
			}
			withTopicsExecuted = true
			return mq.WithTopics(topics...)
		}

		return func() {
			if maxAttempts1 != nil && !withMaxAttemptsExecuted {
				t.Errorf("withMaxAttempts() not executed")
			}
			if topics1 != nil && !withTopicsExecuted {
				t.Errorf("withTopics() not executed")
			}
		}
	}

	type fields struct {
		MaxAttempts int32
		Topics      []string
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
				MaxAttempts: 0,
				Topics:      nil,
			},
		},
		{
			name: "only max attempts",
			setUp: func(fields fields) func() {
				return Expect(&fields.MaxAttempts, nil)
			},
			fields: fields{
				MaxAttempts: 3,
				Topics:      nil,
			},
		},
		{
			name: "only topics",
			setUp: func(fields fields) func() {
				return Expect(nil, fields.Topics)
			},
			fields: fields{
				MaxAttempts: 0,
				Topics:      []string{"topic1", "topic2"},
			},
		},
		{
			name: "both max attempts and topics",
			setUp: func(fields fields) func() {
				return Expect(&fields.MaxAttempts, fields.Topics)
			},
			fields: fields{
				MaxAttempts: 5,
				Topics:      []string{"topic1", "topic2", "topic3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.setUp(tt.fields)()
			p := &ProducerOption{
				MaxAttempts: tt.fields.MaxAttempts,
				Topics:      tt.fields.Topics,
			}
			_ = p.ToOptions()
		})
	}
}

func TestConsumerOption_ToOptions(t *testing.T) {

	var Expect = func(awaitDuration1 *time.Duration, subscriptionExpressions1 map[string]*mq.FilterExpression) func() {
		withAwaitDurationExecuted := false
		withSubscriptionExpressionsExecuted := false

		withAwaitDuration = func(awaitDuration time.Duration) mq.SimpleConsumerOption {
			if *awaitDuration1 != awaitDuration {
				t.Errorf("withAwaitDuration() = %v, want %v", awaitDuration1, awaitDuration)
			}
			withAwaitDurationExecuted = true
			return mq.WithSimpleAwaitDuration(awaitDuration)
		}
		withSubscriptionExpressions = func(subscriptionExpressions map[string]*mq.FilterExpression) mq.SimpleConsumerOption {
			if !reflect.DeepEqual(subscriptionExpressions1, subscriptionExpressions) {
				t.Errorf("withSubscriptionExpressions() = %v, want %v", subscriptionExpressions1, subscriptionExpressions)
			}
			withSubscriptionExpressionsExecuted = true
			return mq.WithSimpleSubscriptionExpressions(subscriptionExpressions)
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
			name: "only topic",
			setUp: func(fields fields) func() {
				return Expect(nil, map[string]*mq.FilterExpression{
					"test-topic": mq.NewFilterExpressionWithType("*", mq.TAG),
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
						Topic: "test-topic",
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

func TestProvideConsumer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockConsumer := NewMockSimpleConsumer(controller)
	mockConsumer.EXPECT().Start().Return(nil)
	mockConsumer.EXPECT().GracefulStop().Return(nil)

	newSimpleConsumer = func(config *mq.Config, opts ...mq.SimpleConsumerOption) (mq.SimpleConsumer, error) {
		return mockConsumer, nil
	}

	gone.
		NewApp(LoadConsumer).
		Run(func(simpleConsumer mq.SimpleConsumer) {
			assert.Equal(t, mockConsumer, simpleConsumer)
		})
}

func TestProvideProducer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockProducer := NewMockProducer(controller)
	mockProducer.EXPECT().Start().Return(nil)
	mockProducer.EXPECT().GracefulStop().Return(nil)

	newProducer = func(config *mq.Config, opts ...mq.ProducerOption) (mq.Producer, error) {
		return mockProducer, nil
	}
	gone.NewApp(LoadProducer).
		Run(func(producer mq.Producer) {
			assert.Equal(t, mockProducer, producer)
		})
}
