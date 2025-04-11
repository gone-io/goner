package cmux

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(s Server) {
			assert.NotNil(t, s)
		})
}
