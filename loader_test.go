package goner

import (
	"github.com/gone-io/gone/v2"
	"testing"
)

func TestGinLoad(t *testing.T) {
	gone.NewApp(BaseLoad).Run()
	gone.NewApp(GinLoad).Run()
}
