module github.com/gone-io/goner/deepseek

go 1.24.1

require (
	github.com/cohesion-org/deepseek-go v1.2.8
	github.com/gone-io/gone/v2 v2.0.10
	github.com/gone-io/goner/g v1.0.8
	github.com/stretchr/testify v1.10.0
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
