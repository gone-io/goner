package urllib

import (
	"github.com/gone-io/gone/v2"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_requestProvider_Provide(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(r *req.Request) {
			assert.NotNil(t, r)
		})
}

func Test_clientProvider_Provide(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(c *req.Client) {
			assert.NotNil(t, c)
		})
}
