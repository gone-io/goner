package apollo

import (
	"github.com/gone-io/gone/v2"
	"testing"
)

func TestLoad(t *testing.T) {
	gone.NewApp(Load).Run()
}
