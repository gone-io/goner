module mcp/quickstart

go 1.24.1

replace github.com/gone-io/goner/mcp => ./../../../mcp

replace github.com/gone-io/goner/g => ../../../g

replace github.com/gone-io/goner/viper => ../../../viper

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/mcp v1.3.6
	github.com/mark3labs/mcp-go v0.39.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/gone-io/goner/g v1.3.6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.uber.org/mock v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
