package xorm

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

func Load(loader gone.Loader) error {
	return load(loader)
}

var load = g.BuildOnceLoadFunc(
	g.L(&xormProvider{}),
	g.L(xormEngineProvider),
	g.L(xormGroupProvider),
	g.L(&engProvider{}, gone.IsDefault(new(Engine), new([]Engine))),
)
