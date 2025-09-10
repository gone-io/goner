module github.com/gone-io/goner/gorm/sqlserver

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/stretchr/testify v1.10.0
	gorm.io/driver/sqlserver v1.6.1
	gorm.io/gorm v1.30.5
)

replace github.com/gone-io/goner/g => ../../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/microsoft/go-mssqldb v1.9.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
