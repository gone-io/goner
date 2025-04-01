package nacos

import "github.com/gone-io/gone/v2"

func Load(loader gone.Loader) error {
	var load = gone.OnceLoad(func(loader gone.Loader) error {
		err := loader.
			Load(
				&configure{},
				gone.Name(gone.ConfigureName),
				gone.IsDefault(new(gone.Configure)),
				gone.ForceReplace(),
			)
		if err != nil {
			return err
		}

		return nil
	})
	return load(loader)
}

func LoadRegistry(loader gone.Loader) error {
	var load = gone.OnceLoad(func(loader gone.Loader) error {
		return loader.Load(&Registry{})
	})
	return load(loader)
}
