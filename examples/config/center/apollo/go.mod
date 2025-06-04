module examples/config/center/appllo

go 1.24.1

require (
	github.com/gone-io/gone/v2 v2.2.5
	github.com/gone-io/goner/apollo v1.3.1
	github.com/gone-io/goner/g v1.3.1
)

replace (
	github.com/gone-io/goner/apollo => ./../../../../apollo
	github.com/gone-io/goner/g => ./../../../../g
	github.com/gone-io/goner/viper => ./../../../../viper
)

require (
	github.com/apolloconfig/agollo/v4 v4.4.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/encoding/javaproperties v0.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gone-io/goner/viper v1.3.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.8.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
