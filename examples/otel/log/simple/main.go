package main

import "github.com/gone-io/gone/v2"

func main() {
	gone.
		Loads(GoneModuleLoad).
		Run(func(logger gone.Logger) {
			logger.Infof("hello world")
			logger.Errorf("error info")
		})
}
