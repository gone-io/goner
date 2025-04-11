module github.com/gone-io/goner/schedule

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.10
	github.com/robfig/cron/v3 v3.0.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.1
)

require github.com/gone-io/goner/g v1.0.9

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
