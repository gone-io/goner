module github.com/gone-io/goner/zap

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.10
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require github.com/gone-io/goner/g v1.0.9

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
