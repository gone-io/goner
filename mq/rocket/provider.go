package rocket

import (
	"fmt"
	mq "github.com/apache/rocketmq-clients/golang"
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

func (o *ConsumerOption) ToOptions() []mq.SimpleConsumerOption {
	options := make([]mq.SimpleConsumerOption, 0)
	if o.AwaitDuration > 0 {
		options = append(options, mq.WithAwaitDuration(o.AwaitDuration))
	}
	if len(o.Expressions) > 0 {
		m := make(map[string]*mq.FilterExpression)
		for _, v := range o.Expressions {
			var t mq.FilterExpressionType
			switch v.Type {
			case "sql":
				t = mq.SQL92
			default:
				t = mq.TAG
			}
			m[v.Topic] = mq.NewFilterExpressionWithType(v.Expression, t)
		}
		options = append(options, mq.WithSubscriptionExpressions(m))
	}
	return options
}

func ProvideConsumer(tagConf string, param struct {
	configure gone.Configure `gone:"configure"`
}) (mq.SimpleConsumer, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("rocketmq.%s", name)

	var conf mq.Config
	err := param.configure.Get(confKey, &conf, "")
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", confKey))
	}

	var option ConsumerOption
	_ = param.configure.Get(fmt.Sprintf("rocketmq.%s.consumer", name), &option, "")

	consumer, err := mq.NewSimpleConsumer(&conf, option.ToOptions()...)
	return consumer, gone.ToErrorWithMsg(err, "can not create rocketmq consumer")
}

type ProducerOption struct {
	MaxAttempts int32    `json:"maxAttempts,omitempty"`
	Topics      []string `json:"topics,omitempty"`
}

func (p *ProducerOption) ToOptions() []mq.ProducerOption {
	options := make([]mq.ProducerOption, 0)
	if p.MaxAttempts > 0 {
		options = append(options, mq.WithMaxAttempts(p.MaxAttempts))
	}
	if len(p.Topics) > 0 {
		options = append(options, mq.WithTopics(p.Topics...))
	}
	return options
}

func ProvideProducer(tagConf string, param struct {
	configure gone.Configure   `gone:"configure"`
	keeper    gone.GonerKeeper `gone:"*"`
}) (mq.Producer, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("rocketmq.%s", name)

	var conf mq.Config
	err := param.configure.Get(confKey, &conf, "")
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", confKey))
	}

	var option ProducerOption
	_ = param.configure.Get(fmt.Sprintf("rocketmq.%s.producer", name), &option, "")
	options := option.ToOptions()

	if checker, err := g.GetComponentByName[*mq.TransactionChecker](param.keeper, fmt.Sprintf("rocketmq.%s", name)); err == nil {
		options = append(options, mq.WithTransactionChecker(checker))
	}

	producer, err := mq.NewProducer(&conf, options...)
	return producer, gone.ToErrorWithMsg(err, "can not create rocketmq producer")
}
