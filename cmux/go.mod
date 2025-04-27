module github.com/gone-io/goner/cmux

go 1.24.1

require (
	github.com/gone-io/gone/mock/v2 v2.1.0
	github.com/gone-io/gone/v2 v2.1.0
	github.com/gone-io/goner/g v1.0.11
	github.com/soheilhy/cmux v0.1.5
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.1
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
