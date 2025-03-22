module grpc_demo

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.7
	github.com/gone-io/goner/grpc v0.0.0-00010101000000-000000000000
	github.com/gone-io/goner/viper v0.0.1
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
)

replace (
	github.com/gone-io/goner/cmux => ../../cmux
	github.com/gone-io/goner/grpc => ../
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/encoding/javaproperties v0.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gone-io/goner/cmux v0.0.0-00010101000000-000000000000 // indirect
	github.com/gone-io/goner/tracer v0.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a // indirect
	github.com/sagikazarmark/locafero v0.8.0 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250311190419-81fb87f6b8bf // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
