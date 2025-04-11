package gin

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"net/http"
)

var load = g.BuildOnceLoadFunc(
	g.L(&router{},
		gone.IsDefault(
			new(RouteGroup),
			new(IRouter),
			new(http.Handler),
		),
	),
	g.L(&SysMiddleware{}),
	g.L(&proxy{}, gone.IsDefault(new(HandleProxyToGin))),
	g.L(NewGinResponser()),
	g.L(&httpInjector{}),
	g.F(func(loader gone.Loader) error {
		return loader.Load(NewGinServer())
	}),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
