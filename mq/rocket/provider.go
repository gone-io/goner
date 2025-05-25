package rocket

import (
	"fmt"
	mq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"time"
)

type ConsumerOption struct {
	AwaitDuration time.Duration `json:"awaitDuration,omitempty"`
	Expressions   []struct {
		Topic      string `json:"topic"`
		Type       string `json:"type"` // sql or tag
		Expression string `json:"expression"`
	} `json:"expressions"`
}

var withAwaitDuration = mq.WithAwaitDuration
var withSubscriptionExpressions = mq.WithSubscriptionExpressions

func (o *ConsumerOption) ToOptions() []mq.SimpleConsumerOption {
	options := make([]mq.SimpleConsumerOption, 0)
	if o.AwaitDuration > 0 {
		options = append(options, withAwaitDuration(o.AwaitDuration))
	}
	if len(o.Expressions) > 0 {
		m := make(map[string]*mq.FilterExpression)
		for _, v := range o.Expressions {
			var t mq.FilterExpressionType
			switch v.Type {
			case "sql":
				t = mq.SQL92
				m[v.Topic] = mq.NewFilterExpressionWithType(v.Expression, t)
			case "tag":
				t = mq.TAG
				m[v.Topic] = mq.NewFilterExpressionWithType(v.Expression, t)
			default:
				m[v.Topic] = mq.SUB_ALL
			}
		}
		options = append(options, withSubscriptionExpressions(m))
	}
	return options
}

var newSimpleConsumer = mq.NewSimpleConsumer

func ProvideConsumer(tagConf string, param struct {
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (mq.SimpleConsumer, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("rocketmq.%s", name)

	var conf mq.Config
	err := param.configure.Get(confKey, &conf, "")
	g.PanicIfErr(gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", confKey)))

	var option ConsumerOption
	_ = param.configure.Get(fmt.Sprintf("rocketmq.%s.consumer", name), &option, "")

	consumer, err := newSimpleConsumer(&conf, option.ToOptions()...)

	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not create rocketmq consumer"))

	err = consumer.Start()
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not start rocketmq consumer"))

	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, consumer.GracefulStop(), "close kafka consumer group failed")
	})

	return consumer, nil
}

type ProducerOption struct {
	MaxAttempts int32    `json:"maxAttempts,omitempty"`
	Topics      []string `json:"topics,omitempty"`
}

var withMaxAttempts = mq.WithMaxAttempts
var withTopics = mq.WithTopics

func (p *ProducerOption) ToOptions() []mq.ProducerOption {
	options := make([]mq.ProducerOption, 0)
	if p.MaxAttempts > 0 {
		options = append(options, withMaxAttempts(p.MaxAttempts))
	}
	if len(p.Topics) > 0 {
		options = append(options, withTopics(p.Topics...))
	}
	return options
}

var newProducer = mq.NewProducer

func ProvideProducer(tagConf string, param struct {
	configure  gone.Configure   `gone:"configure"`
	keeper     gone.GonerKeeper `gone:"*"`
	beforeStop gone.BeforeStop  `gone:"*"`
	logger     gone.Logger      `gone:"*"`
}) (mq.Producer, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("rocketmq.%s", name)

	var conf mq.Config
	err := param.configure.Get(confKey, &conf, "")
	g.PanicIfErr(gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", confKey)))
	var option ProducerOption
	_ = param.configure.Get(fmt.Sprintf("rocketmq.%s.producer", name), &option, "")
	options := option.ToOptions()

	if checker, err := g.GetComponentByName[*mq.TransactionChecker](param.keeper, fmt.Sprintf("rocketmq.%s", name)); err == nil {
		options = append(options, mq.WithTransactionChecker(checker))
	}

	producer, err := newProducer(&conf, options...)
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not create rocketmq producer"))

	err = producer.Start()
	g.PanicIfErr(gone.ToErrorWithMsg(err, "can not start rocketmq producer"))

	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, producer.GracefulStop(), "close kafka consumer group failed")
	})
	return producer, nil
}
