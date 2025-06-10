module mcp/quickstart

go 1.24.1

replace github.com/gone-io/goner/mcp => ./../../../mcp

replace github.com/gone-io/goner/g => ../../../g

replace github.com/gone-io/goner/viper => ../../../viper

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/mcp v1.3.2
	github.com/mark3labs/mcp-go v0.23.1
)

require (
	github.com/gone-io/goner/g v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.uber.org/mock v0.5.2 // indirect
)
