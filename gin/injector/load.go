package injector

import "github.com/gone-io/gone/v2"

func BuildLoad[T any](name string) gone.LoadFunc {
	return func(loader gone.Loader) error {
		loader.
			MustLoad(&delayBindInjector[T]{name: name}).
			MustLoad(&bindExecutor[T]{})
		return nil
	}
}
