package nacos

import (
	"github.com/gone-io/gone/v2"
	viper "github.com/gone-io/goner/viper"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"testing"
)

func Test_configure_init(t *testing.T) {
	gone.
		NewApp(viper.Load).
		Test(func(localConf gone.Configure) {
			c := configure{}
			init, err := c.init(localConf, clients.NewConfigClient)
			if err != nil {
				t.Error(err)
			}
			if init == nil {
				t.Error("init is nil")
			}
		})
}

func Test_configure_Init(t *testing.T) {
	t.Skip("Skip integration test that requires Nacos server")

	gone.
		NewApp(viper.Load).
		Loads(Load).
		Test(func() {
			//todo
		})
}
