package apollo

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var load = g.BuildOnceLoadFunc(
	g.L(&apolloConfigure{},
		gone.Name(gone.ConfigureName),
		gone.IsDefault(new(gone.Configure)),
		gone.ForceReplace(),
	),
	g.L(&changeListener{}),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
