module examples/config_center/appllo

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.0.10
	github.com/gone-io/goner/apollo v1.0.2
)

replace github.com/gone-io/goner/viper v1.0.4 => ../../../viper

replace github.com/gone-io/goner/apollo => ../../../apollo

require (
	github.com/apolloconfig/agollo/v4 v4.4.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/encoding/javaproperties v0.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gone-io/goner/viper v1.0.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
