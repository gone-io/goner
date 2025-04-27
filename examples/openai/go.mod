module github.com/gone-io/goner/examples/openai

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.1.0
	github.com/gone-io/goner/openai v1.1.0
	github.com/sashabaranov/go-openai v1.38.1
)

require github.com/gone-io/goner/g v1.1.0 // indirect

replace github.com/gone-io/goner/openai => ../../openai

replace github.com/gone-io/goner/g => ../../g
