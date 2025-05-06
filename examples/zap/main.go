package main

import "github.com/gone-io/gone/v2"

func main() {
	gone.
		Loads(GoneModuleLoad).
		Load(&UseOriginZap{}).
		Load(&UseGoneLogger{}).
		Load(&UseTracer{}).
		Run(func(p1 *UseOriginZap, p2 *UseGoneLogger, p3 *UseTracer) {
			p1.PrintLog()
			p2.PrintLog()
			p3.PrintLog()
		})
}
