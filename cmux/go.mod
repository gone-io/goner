module github.com/gone-io/goner/cmux

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.7
	github.com/gone-io/goner/g v0.0.0-00010101000000-000000000000
	github.com/soheilhy/cmux v0.1.5
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.0
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb // indirect
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
