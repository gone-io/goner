package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gone-io/gone/v2"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Skip("skip integration test")
	// cd testdata && docker compose up -d

	_ = os.Setenv("GONE_MQTT_DEFAULT", `{
"brokers": ["localhost:1883"],
"ClientID": "ClientID-test"
}`)
	defer func() {
		_ = os.Unsetenv("GONE_MQTT_DEFAULT")
	}()

	gone.
		NewApp(Load).
		Run(func(client mqtt.Client) {
			topic := "topic/test"
			info := "hello gone"

			ch := make(chan struct{})

			token := client.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
				fmt.Printf("%s", message)
				if string(message.Payload()) == info {
					close(ch)
				}
			})
			token.Wait()
			if err := token.Error(); err != nil {
				t.Errorf("subscribe error: %s", err)
			}

			token = client.Publish(topic, 1, false, info)
			token.Wait()
			token.Wait()
			if err := token.Error(); err != nil {
				t.Errorf("publish error: %s", err)
			}

			<-ch
		})
}
