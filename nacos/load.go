package nacos

import "github.com/gone-io/gone/v2"

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

func Load(loader gone.Loader) error {
	return load(loader)
}
