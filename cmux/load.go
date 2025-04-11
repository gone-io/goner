package cmux

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"net"
)

var load = g.BuildOnceLoadFunc(
	g.L(
		&server{listen: net.Listen},
		gone.IsDefault(new(CMuxServer)),
		gone.HighStartPriority(),
	),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
