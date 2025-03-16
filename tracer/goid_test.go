package tracer

import (
	"github.com/petermattis/goid"
	"testing"
)

func Test_GetGoId(t *testing.T) {
	gid := goid.Get()
	if gid == 0 {
		t.Fatal("can not get goid")
	}
	t.Logf("gid=%d", gid)
}
