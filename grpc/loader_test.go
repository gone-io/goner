package grpc

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	_ = os.Setenv("GONE_SERVER_GRPC_PORT", "0")
	defer func() {
		_ = os.Unsetenv("GONE_SERVER_GRPC_PORT")
	}()

	gone.
		NewApp(Load).
		Run(func(s *server, c *clientRegister) {
			assert.NotNil(t, s)
			assert.NotNil(t, c)
		})

	gone.
		NewApp(ServerLoad).
		Run(func(s *server) {
			assert.NotNil(t, s)
		})

	gone.
		NewApp(ClientLoad).
		Run(func(c *clientRegister) {
			assert.NotNil(t, c)
		})
}
