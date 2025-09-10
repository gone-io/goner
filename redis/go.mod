module github.com/gone-io/goner/redis

go 1.24.1

require (
	github.com/bytedance/sonic v1.14.1
	github.com/goccy/go-json v0.10.5
	github.com/gomodule/redigo v1.9.2
	github.com/gone-io/gone/v2 v2.2.6
	github.com/google/uuid v1.6.0
	github.com/json-iterator/go v1.1.12
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.6.0
)

require github.com/gone-io/goner/g v1.3.5

replace github.com/gone-io/goner/g => ../g

require (
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.21.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
