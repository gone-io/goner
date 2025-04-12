package urllib

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var load = g.BuildOnceLoadFunc(
	g.L(&r{}, gone.IsDefault(new(Client))),
	g.L(&requestProvider{}),
	g.L(&clientProvider{}),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
