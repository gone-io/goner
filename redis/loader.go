package redis

import (
	"github.com/gone-io/gone/v2"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&inner{}).
		MustLoad(&pool{}, gone.IsDefault(new(Pool))).
		MustLoad(&cache{}, gone.IsDefault(new(Cache), new(Key))).
		MustLoad(&locker{}, gone.IsDefault(new(Locker))).
		MustLoad(&provider{}, gone.IsDefault(new(HashProvider)))
	return nil
}
