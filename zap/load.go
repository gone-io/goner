package gone_zap

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

func Load(loader gone.Loader) error {
	//defaultLog, _ := zap.NewDevelopment()
	loader.
		MustLoad(&atomicLevel{}).
		MustLoad(&zapLoggerProvider{}).
		MustLoad(&ctxLogger{}, gone.IsDefault(new(g.CtxLogger))).
		MustLoad(
			&goneLogger{},
			gone.IsDefault(new(gone.Logger)),
			gone.ForceReplace(),
		)
	return nil
}
