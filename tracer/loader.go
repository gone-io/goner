package tracer

import "github.com/gone-io/gone/v2"

func Load(loader gone.Loader) error {
	return loader.Load(
		&tracer{},
		gone.IsDefault(new(Tracer)),
	)
}

func LoadGidTracer(loader gone.Loader) error {
	return loader.Load(
		&tracerOverGid{},
		gone.IsDefault(new(Tracer)),
	)
}
