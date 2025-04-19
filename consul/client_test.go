package consul

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
)

func TestClientLoad(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		_ = os.Setenv("GONE_CONSUL", `{"address":"x://127.0.0.", "Scheme":"x"}`)
		defer func() {
			_ = os.Unsetenv("GONE_CONSUL")
		}()

		gone.
			NewApp(ClientLoad).
			Run(func(provider gone.Provider[*api.Client]) {
				provide, err := provider.Provide("")
				assert.Nil(t, provide)
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "can not create consul client")
			})
	})

	t.Run("success", func(t *testing.T) {
		gone.
			NewApp(ClientLoad).
			Run(func(client *api.Client, in struct {
				c2 *api.Client `gone:"*"`
			}) {
				assert.Equal(t, client, in.c2)

				k := "test"
				v := fmt.Sprintf("test-%d", rand.Int())

				put, err := client.KV().Put(&api.KVPair{Key: k,
					Value: []byte(v),
				}, nil)
				assert.Nil(t, err)
				assert.NotNil(t, put)

				get, _, err := client.KV().Get(k, nil)
				assert.Nil(t, err)
				assert.NotNil(t, get)
				assert.Equal(t, v, string(get.Value))
			})
	})
}
