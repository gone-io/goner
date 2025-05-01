package gone_zap

import (
	"github.com/gone-io/gone/v2"
)

func Load(loader gone.Loader) error {
	//defaultLog, _ := zap.NewDevelopment()
	loader.
		MustLoad(&atomicLevel{}).
		MustLoad(&zapLoggerProvider{}).
		MustLoad(&sugarProvider{}).
		MustLoad(
			&sugar{
				//SugaredLogger: defaultLog.Sugar(),
			},
			gone.IsDefault(new(gone.Logger)),
			gone.ForceReplace(),
		)
	return nil
}
