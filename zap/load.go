package gone_zap

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var load = g.BuildOnceLoadFunc(
	g.L(&atomicLevel{}),
	g.L(&zapLoggerProvider{}),
	g.L(&sugarProvider{}),
	g.L(&sugar{}, gone.IsDefault(new(gone.Logger)), gone.ForceReplace()),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
