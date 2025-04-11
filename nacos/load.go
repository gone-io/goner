package nacos

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var load = g.BuildOnceLoadFunc(g.L(&configure{},
	gone.Name(gone.ConfigureName),
	gone.IsDefault(new(gone.Configure)),
	gone.ForceReplace(),
))

func Load(loader gone.Loader) error {
	return load(loader)
}

var registryLoad = g.BuildOnceLoadFunc(g.L(&Registry{}))

func RegistryLoad(loader gone.Loader) error {
	return registryLoad(loader)
}
