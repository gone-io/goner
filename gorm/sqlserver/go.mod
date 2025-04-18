module github.com/gone-io/goner/gorm/sqlserver

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.11
	github.com/stretchr/testify v1.10.0
	gorm.io/driver/sqlserver v1.5.4
	gorm.io/gorm v1.25.12
)

require github.com/gone-io/goner/g v1.0.10

replace github.com/gone-io/goner/g => ../../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/microsoft/go-mssqldb v1.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
