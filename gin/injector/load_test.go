package injector

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildLoad(t *testing.T) {
	gone.
		NewApp(BuildLoad[any]("test")).
		Run(func(i DelayBindInjector[any]) {
			assert.NotNil(t, i)
		})
}
