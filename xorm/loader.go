package xorm

import (
	"github.com/gone-io/gone/v2"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&xormProvider{}).
		MustLoad(xormEngineProvider).
		MustLoad(xormGroupProvider).
		MustLoad(&engProvider{}, gone.IsDefault(new(Engine), new([]Engine)))
	return nil
}
