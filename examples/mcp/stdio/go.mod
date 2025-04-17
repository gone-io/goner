module mcp/stdio

go 1.24.1

replace github.com/gone-io/goner/mcp => ./../../../mcp

replace github.com/gone-io/goner/g => ../../../g

replace github.com/gone-io/goner/viper => ../../../viper

require (
	github.com/gone-io/gone/v2 v2.0.11
	github.com/gone-io/goner/mcp v1.0.9
	github.com/gone-io/goner/viper v1.0.9
	github.com/mark3labs/mcp-go v0.21.0
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/encoding/javaproperties v0.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gone-io/goner/g v1.0.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
