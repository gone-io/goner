module github.com/gone-io/goner/gin

go 1.24.1

require (
	github.com/bytedance/sonic v1.13.2
	github.com/gin-gonic/gin v1.10.0
	github.com/goccy/go-json v0.10.5
	github.com/gone-io/gone/v2 v2.0.6
	github.com/gone-io/goner/cmux v0.0.0-00010101000000-000000000000
	github.com/gone-io/goner/tracer v0.0.1
	github.com/google/uuid v1.6.0
	github.com/json-iterator/go v1.1.12
	github.com/pkg/errors v0.9.1
	github.com/soheilhy/cmux v0.1.5
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.0
	golang.org/x/time v0.11.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.0.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.25.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/gone-io/goner/cmux => ../cmux
	github.com/gone-io/goner/tracer => ../tracer
)
