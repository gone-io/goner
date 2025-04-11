package cmux

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"net"
)

func Load(loader gone.Loader) error {
	return g.BuildLoadFunc(loader, g.L(
		&server{listen: net.Listen},
		gone.IsDefault(new(CMuxServer)),
		gone.HighStartPriority()),
	)
}
