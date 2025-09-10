module examples/openai

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/openai v1.3.5
	github.com/sashabaranov/go-openai v1.41.1
)

require (
	github.com/gone-io/goner/g v1.3.5 // indirect
	go.uber.org/mock v0.6.0 // indirect
)

replace github.com/gone-io/goner/openai => ../../openai

replace github.com/gone-io/goner/g => ../../g
