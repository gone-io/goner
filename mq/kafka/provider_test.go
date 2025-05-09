package kafka

import (
	"github.com/IBM/sarama"
	"github.com/gone-io/gone/mock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestConf_ReadFromConfigure(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	configure := mock.NewMockConfigure(controller)
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
