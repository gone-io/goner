package provide_func

import "github.com/gone-io/gone/v2"

// ThirdComponent Emulating third-party component
type ThirdComponent struct {
}

func provide(tagConf string, in struct {
	logger gone.Logger `gone:"*"`
}) (*ThirdComponent, error) {
	confMap, confKeys := gone.TagStringParse(tagConf)
	in.logger.Infof("confMap => %#v\nconfKeys=>%#v", confMap, confKeys)

	// create third-party component for different conf
	return &ThirdComponent{}, nil
}

// Load dine gone.LoadFunc for ThirdComponent
func Load(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(provide))
}
