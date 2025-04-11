package urllib

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

func Load(loader gone.Loader) error {
	return g.BuildLoadFunc(loader,
		g.L(&r{}, gone.IsDefault(new(Client))),
		g.L(&requestProvider{}),
		g.L(&clientProvider{}),
	)
}
