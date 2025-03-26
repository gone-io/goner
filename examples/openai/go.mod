module github.com/gone-io/goner/examples/openai

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.9
	github.com/gone-io/goner/openai v0.0.0
	github.com/sashabaranov/go-openai v1.38.1
)

replace github.com/gone-io/goner/openai => ../../openai
