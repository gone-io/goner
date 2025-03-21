module github.com/gone-io/goner/gorm/grpc

go 1.24.1

replace (
	github.com/gone-io/goner/cmux => ../cmux
	github.com/gone-io/goner/tracer => ../tracer
)

require (
	github.com/gone-io/gone/v2 v2.0.6
	github.com/gone-io/goner/cmux v0.0.0-00010101000000-000000000000
	github.com/gone-io/goner/tracer v0.0.1
	github.com/soheilhy/cmux v0.1.5
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.0
	google.golang.org/grpc v1.71.0
	xorm.io/xorm v1.3.9
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goccy/go-json v0.8.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	xorm.io/builder v0.3.11-0.20220531020008-1bd24a7dc978 // indirect
)
