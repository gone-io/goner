module github.com/gone-io/goner/xorm

go 1.24.1

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/gone-io/gone/mock/v2 v2.0.13
	github.com/gone-io/gone/v2 v2.0.14
	github.com/gone-io/goner/g v1.0.11
	github.com/jmoiron/sqlx v1.4.0
	github.com/jtolds/gls v4.20.0+incompatible
	github.com/spf13/cast v1.7.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.1
	xorm.io/xorm v1.3.9
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	xorm.io/builder v0.3.13 // indirect
)
