package gin

import (
	"github.com/gone-io/gone/v2"
	"net/http"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&router{},
			gone.IsDefault(
				new(RouteGroup),
				new(IRouter),
				new(http.Handler),
			),
		).
		MustLoad(&SysMiddleware{}).
		MustLoad(&proxy{}, gone.IsDefault(new(HandleProxyToGin))).
		MustLoad(NewGinResponser()).
		MustLoad(&httpInjector{})
	return loader.Load(NewGinServer())
}
