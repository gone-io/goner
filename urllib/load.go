package urllib

import (
	"github.com/gone-io/gone/v2"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&r{}, gone.IsDefault(new(Client))).
		MustLoad(&requestProvider{}).
		MustLoad(&clientProvider{})
	return nil
}
