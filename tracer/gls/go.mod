module github.com/gone-io/goner/tracer/gls

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.2
	github.com/gone-io/goner/g v1.2.1-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/jtolds/gls v4.20.0+incompatible
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/mock v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gone-io/goner/g => ../../g
