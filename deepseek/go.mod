module github.com/gone-io/goner/deepseek

go 1.24.1

require (
	github.com/cohesion-org/deepseek-go v1.3.2
	github.com/gone-io/gone/v2 v2.2.6
	github.com/stretchr/testify v1.10.0
)

require github.com/gone-io/goner/g v1.3.5

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/ollama/ollama v0.11.10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
