package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/balancer"
	"github.com/gone-io/goner/nacos"
	"github.com/gone-io/goner/urllib"
	"github.com/gone-io/goner/viper"
	"time"
)

func main() {
	gone.
		NewApp(
			nacos.RegistryLoad,
			balancer.Load,
			viper.Load,
			urllib.Load,
		).
		Run(func(client urllib.Client, logger gone.Logger) {

			call := func() {
				var data urllib.Res[string]
				res, err := client.
					R().
					SetSuccessResult(&data).
					Get("http://user-center/hello?name=goner")
				if err != nil {
					logger.Errorf("client request err: %v", err)
					return
				}

				if res.IsSuccessState() {
					logger.Infof("res=> %#v", data)
				}
			}

			for i := 0; i < 10; i++ {
				call()
			}

			time.Sleep(10 * time.Second)
		})
}
