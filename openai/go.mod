module github.com/gone-io/goner/openai

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.2
	github.com/sashabaranov/go-openai v1.38.1
	github.com/stretchr/testify v1.10.0
)

require github.com/gone-io/goner/g v1.3.0

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/mock v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
