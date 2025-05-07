module github.com/gone-io/goner/openai

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.1.0
	github.com/sashabaranov/go-openai v1.38.1
	github.com/stretchr/testify v1.10.0
)

require github.com/gone-io/goner/g v1.2.1

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
