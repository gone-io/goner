package main

import (
	"github.com/gone-io/gone/v2"
	mcp "github.com/gone-io/goner/mcp"
	"github.com/gone-io/goner/viper"
)

//go:generate gonectr generate
func main() {
	// IF do not use viper
	//_ = os.Setenv("GONE_MCP", `{"name":"demo", "version":"1.0.0", "transportType":"sse"}`)
	//_ = os.Setenv("GONE_MCP_SSE", `{"address":":8082"}`)
	//defer func() {
	//	_ = os.Unsetenv("GONE_MCP")
	//	_ = os.Unsetenv("GONE_MCP_SSE")
	//}()

	gone.
		Loads(mcp.ServerLoad, viper.Load).
		Serve()
}
