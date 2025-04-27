package schedule

import "github.com/gone-io/gone/v2"

func Load(loader gone.Loader) error {
	return loader.Load(&schedule{})
}
