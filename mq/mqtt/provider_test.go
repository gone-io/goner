package mqtt

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestProvideClient(t *testing.T) {
	_ = os.Setenv("GONE_MQTT_DEFAULT", `{
"brokers": ["localhost:8883"]
}`)
	defer func() {
		_ = os.Unsetenv("GONE_MQTT_DEFAULT")
	}()

	controller := gomock.NewController(t)
	defer controller.Finish()
	client := NewMockClient(controller)

	t.Run("connect success", func(t *testing.T) {
		newClient = func(o *mqtt.ClientOptions) mqtt.Client {
			assert.Len(t, o.Servers, 1)
			assert.Equal(t, "tcp://localhost:8883", o.Servers[0].String())

			token := NewMockToken(controller)

			client.EXPECT().Connect().Return(token)
			token.EXPECT().Wait().Return(true)
			token.EXPECT().Error().Return(nil)
			client.EXPECT().Disconnect(gomock.Any())
			return client
		}

		gone.
			NewApp(Load).
			Run(func(client mqtt.Client) {
				assert.NotNil(t, client)
			})
	})

	t.Run("connect error", func(t *testing.T) {
		newClient = func(o *mqtt.ClientOptions) mqtt.Client {
			assert.Len(t, o.Servers, 1)
			assert.Equal(t, "tcp://localhost:8883", o.Servers[0].String())

			token := NewMockToken(controller)

			client.EXPECT().Connect().Return(token)
			token.EXPECT().Wait().Return(true)
			token.EXPECT().Error().Return(errors.New("connect error"))

			return client
		}

		err := gone.SafeExecute(func() error {
			gone.
				NewApp(Load).
				Run(func(client mqtt.Client) {
					assert.NotNil(t, client)
				})

			return nil
		})

		assert.Error(t, err)
	})
}
