package cmux

import (
	"github.com/gone-io/gone/v2"
	"net"
)

func Load(loader gone.Loader) error {
	return loader.Load(
		&server{listen: net.Listen},
		gone.IsDefault(new(CMuxServer)),
		gone.HighStartPriority(),
	)
}
