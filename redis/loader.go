package redis

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

var load = g.BuildOnceLoadFunc(
	g.F(func(loader gone.Loader) error {
		return loader.Load(&pool{}, gone.IsDefault(new(Pool)))
	}),
	g.L(&inner{}),
	g.L(&cache{}, gone.IsDefault(new(Cache), new(Key))),
	g.L(&locker{}, gone.IsDefault(new(Locker))),
	g.L(&provider{}, gone.IsDefault(new(HashProvider))),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
