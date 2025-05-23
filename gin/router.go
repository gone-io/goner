package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
)

var incr = 0

type router struct {
	gone.Flag

	id              int
	r               gin.IRouter
	*gin.Engine     `gone:"*" option:"allowNil"`
	isOtelLogLoaded g.IsOtelTracerLoaded `gone:"*" option:"allowNil"`
	logger          gone.Logger          `gone:"*"`
	middlewares     []Middleware         `gone:"*"`

	htmlTpl     string `gone:"config,server.html-tpl-pattern"`
	mode        string `gone:"config,server.mode,default=release"`
	serviceName string `gone:"config,server.service-name=gin"`

	HandleProxyToGin `gone:"gone-gin-proxy"`
}

type logWriter struct {
	write func(p []byte) (n int, err error)
}

func (w logWriter) Write(p []byte) (n int, err error) {
	return w.write(p)
}

func debugWriter(logger gone.Logger) logWriter {
	return logWriter{
		write: func(p []byte) (n int, err error) {
			logger.Debugf("%s", p)
			return len(p), nil
		},
	}
}
func errorWriter(logger gone.Logger) logWriter {
	return logWriter{
		write: func(p []byte) (n int, err error) {
			logger.Errorf("%s", p)
			return len(p), nil
		},
	}
}

func (r *router) getMiddlewaresFunc() (list []gin.HandlerFunc) {
	for _, middleware := range r.middlewares {
		list = append(list, middleware.Process)
	}
	return list
}

func (r *router) GonerName() string {
	return IdGoneGinRouter
}

func (r *router) Init() {
	if r.Engine == nil {
		gin.SetMode(r.mode)
		r.Engine = gin.New()
	}

	if r.isOtelLogLoaded {
		r.Engine.Use(otelgin.Middleware(r.serviceName))
	}

	r.Engine.Use(r.getMiddlewaresFunc()...)

	if r.htmlTpl != "" {
		r.Engine.LoadHTMLGlob(r.htmlTpl)
	}

	gin.DefaultWriter = debugWriter(r.logger)
	gin.DefaultErrorWriter = errorWriter(r.logger)
}

func (r *router) GetGinRouter() gin.IRouter {
	return r.Engine
}

func (r *router) getR() gin.IRouter {
	if r.r == nil {
		r.r = r.Engine
	}
	return r.r
}

func (r *router) Use(middleware ...HandlerFunc) IRoutes {
	r.getR().Use(r.ProxyForMiddleware(middleware...)...)
	return r
}

func (r *router) Group(relativePath string, handlers ...HandlerFunc) RouteGroup {
	incr++
	return &router{
		id:               incr,
		r:                r.getR().Group(relativePath, r.ProxyForMiddleware(handlers...)...),
		Engine:           r.Engine,
		HandleProxyToGin: r.HandleProxyToGin,
	}
}

func (r *router) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(httpMethod, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodGet, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodPost, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodPut, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodPatch, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodHead, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodOptions, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodDelete, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodConnect, relativePath, r.Proxy(handlers...)...)
	r.getR().Handle(http.MethodTrace, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodGet, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodPost, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodDelete, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodPatch, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodPut, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodOptions, relativePath, r.Proxy(handlers...)...)
	return r
}
func (r *router) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.getR().Handle(http.MethodHead, relativePath, r.Proxy(handlers...)...)
	return r
}
