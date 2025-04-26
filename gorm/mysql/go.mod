module github.com/gone-io/goner/gorm/mysql

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.14
	github.com/stretchr/testify v1.10.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

require github.com/gone-io/goner/g v1.0.11

replace github.com/gone-io/goner/g => ../../g

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.9.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
