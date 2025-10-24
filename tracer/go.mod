module github.com/gone-io/goner/tracer

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/g v1.3.6
	github.com/gone-io/goner/tracer/gid v1.3.6-00010101000000-000000000000
	github.com/gone-io/goner/tracer/gls v1.3.6-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/petermattis/goid v0.0.0-20250904145737-900bdf8bb490 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/mock v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gone-io/goner/g => ../g
	github.com/gone-io/goner/tracer/gid => ./gid
	github.com/gone-io/goner/tracer/gls => ./gls
)
