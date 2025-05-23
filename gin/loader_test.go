package gin

import (
	"github.com/gone-io/gone/v2"
	"testing"
)

func Test_Load(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func() {})
}
