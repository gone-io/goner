package consul

import (
	"testing"
	"time"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
)

func TestRegistryLoad(t *testing.T) {
	t.Run("all function success", func(t *testing.T) {
		gone.
			NewApp(RegistryLoad).
			Run(func(r *Registry) {

				serviceName := "x-test.svc"

				service1 := g.NewService(serviceName, "10.0.11.1", 200, nil, true, 40)
				service2 := g.NewService(serviceName, "10.0.11.2", 200, nil, true, 40)

				err := r.Register(service1)
				assert.Nil(t, err)
				err = r.Register(service2)
				assert.Nil(t, err)

				instances, err := r.GetInstances(serviceName)
				assert.Nil(t, err)
				assert.Equal(t, 2, len(instances))

				ch, stop, err := r.Watch(serviceName)
				assert.Nil(t, err)
				defer func() {
					assert.Nil(t, stop())
				}()
				go func() {
					time.Sleep(1 * time.Second)
					var err = r.Deregister(service1)
					assert.Nil(t, err)
				}()

				services := <-ch
				services = <-ch
				assert.Equal(t, 1, len(services))
				assert.Equal(t, float64(40), services[0].GetWeight())
				assert.Equal(t, "10.0.11.2", services[0].GetIP())

				err = r.Deregister(service2)
				assert.Nil(t, err)
			})
	})
}
