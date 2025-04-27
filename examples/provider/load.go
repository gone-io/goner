package provider

import "github.com/gone-io/gone/v2"

// Load dine gone.LoadFunc for ThirdComponent
func Load(loader gone.Loader) error {
	loader.
		MustLoad(&provider{}).
		MustLoad(&noneParamProvider{})
	return nil
}
