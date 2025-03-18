package apollo

import (
	"github.com/gone-io/gone/v2"
	"testing"
	"time"
)

func Test(t *testing.T) {
	gone.
		NewApp(Load).
		Test(func(in struct {
			a string  `gone:"config,test.a"`
			b *string `gone:"config,test.b"`
			c string  `gone:"config,test.c"`
			d string  `gone:"config,test.d"`
			e string  `gone:"config,test.e"`
			x *string `gone:"config,test.x"`
		}) {
			for true {
				t.Logf("a=%s, b=%s, c=%s, d=%s, e=%s, x=%s", in.a, *in.b, in.c, in.d, in.e, *in.x)
				time.Sleep(10 * time.Second)
			}
		})
}

func Test2(t *testing.T) {
	m := make(map[string]any)

	x := func(x map[string]any) {
		x["test"] = "test"
	}
	x(m)

	t.Log(m)
}
