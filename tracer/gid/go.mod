module github.com/gone-io/goner/tracer/gid

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.1.0
	github.com/gone-io/goner/g v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gone-io/goner/g => ../../g
