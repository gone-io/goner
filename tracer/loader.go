package tracer

import "github.com/gone-io/gone/v2"

var load = gone.OnceLoad(func(loader gone.Loader) error {
	return loader.Load(
		&tracer{},
		gone.IsDefault(new(Tracer)),
		gone.LazyFill(),
	)
})

func Load(loader gone.Loader) error {
	return load(loader)
}

var gidTracerLoad = gone.OnceLoad(func(loader gone.Loader) error {
	return loader.Load(
		&tracerOverGid{},
		gone.IsDefault(new(Tracer)),
		gone.LazyFill(),
	)
})

func LoadGidTracer(loader gone.Loader) error {
	return gidTracerLoad(loader)
}
