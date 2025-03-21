module github.com/gone-io/goner/schedule

go 1.24.1

replace (
	github.com/gone-io/goner/redis => ../redis
	github.com/gone-io/goner/tracer => ../tracer
)

require (
	github.com/gone-io/gone/v2 v2.0.7
	github.com/gone-io/goner/redis v0.0.0-00010101000000-000000000000
	github.com/gone-io/goner/tracer v0.0.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.0
)

require (
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
