module github.com/gone-io/goner/mcp

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.1.0
	github.com/gone-io/goner/g v1.1.1
	github.com/mark3labs/mcp-go v0.23.1
	github.com/stretchr/testify v1.10.0
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
