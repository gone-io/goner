module github.com/gone-io/goner/gorm

go 1.24.1

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.2
	gorm.io/gorm v1.25.12
)

require (
	github.com/gone-io/gone/v2 v2.2.2
	github.com/gone-io/goner/g v1.3.0
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
