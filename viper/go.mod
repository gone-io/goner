module github.com/gone-io/goner/viper

go 1.24.1

require (
	github.com/go-viper/encoding/javaproperties v0.1.0
	github.com/gone-io/gone/v2 v2.2.6
	github.com/gone-io/goner/g v1.3.2
	github.com/sagikazarmark/locafero v0.9.0
	github.com/spf13/afero v1.14.0
	github.com/spf13/viper v1.20.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.2
)

replace github.com/gone-io/goner/g => ../g

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
