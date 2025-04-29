package main

import (
	"github.com/gone-io/gone/v2"
)

//go:generate gonectr generate -m . -s ..
func main() {
	gone.Serve()
}
