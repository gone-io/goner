package gin

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"net/http"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&router{},
			gone.IsDefault(
				new(RouteGroup),
				new(IRouter),
				new(http.Handler),
				new(g.IRoutes),
			),
		).
		MustLoad(&SysMiddleware{}).
		MustLoad(&proxy{}, gone.IsDefault(new(HandleProxyToGin))).
		MustLoad(NewGinResponser()).
		MustLoad(&httpInjector{})
	return loader.Load(NewGinServer())
}
