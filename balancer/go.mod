module github.com/gone-io/goner/balancer

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.1
	github.com/gone-io/goner/g v1.2.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.2
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gone-io/goner/g => ../g
