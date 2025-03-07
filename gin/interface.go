package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"reflect"
)

//go:generate sh -c "mockgen -package=gin github.com/gin-gonic/gin ResponseWriter > gin_response_writer_mock_test.go"
//go:generate sh -c "mockgen -package=gin net Listener > net_listener_mock_test.go"
//go:generate sh -c "mockgen -package=gin -source=../../gin_interface.go |gone mock -o gone_gin_mock_test.go"
//go:generate sh -c "mockgen -package=gin -source=../../interface.go |gone mock -o gone_mock_test.go"
//go:generate sh -c "mockgen -package=gin -self_package=github.com/gone-io/goner/gin -source=interface.go |gone mock -o mock_test.go"

// RouterGroupName Router group name
type RouterGroupName string

type OriginContent = gin.Context

type MountError = GinMountError

// Controller interface, implemented by business code, used to mount and handle routes
// For usage reference [example code](https://gitlab.openviewtech.com/gone/gone-example/-/tree/master/gone-app)
type Controller interface {
	// Mount Route mount interface, this interface will be called before the service starts, the implementation of this function should usually return `nil`
	Mount() MountError
}

// HandleProxyToGin Proxy, provides a proxy function to convert `gone.HandlerFunc` to `gin.HandlerFunc`
// Inject `gin.HandleProxyToGin` using Id: sys-gone-proxy (`gin.SystemGoneProxy`)
type HandleProxyToGin interface {
	Proxy(handler ...HandlerFunc) []gin.HandlerFunc
	ProxyForMiddleware(handlers ...HandlerFunc) (arr []gin.HandlerFunc)
}

type XContext interface {
	JSON(code int, obj any)
	String(code int, format string, values ...any)
	Abort()
}

type WrappedDataFunc func(code int, msg string, data any) any

type WrappedDataFuncSetter interface {
	SetWrappedDataFunc(wrappedDataFunc WrappedDataFunc)
}

// Responser Response handler
// Inject default response handler using Id: gone-gin-responser (`gone.IdGoneGinResponser`)
type Responser interface {
	Success(ctx XContext, data any)
	Failed(ctx XContext, err error)
	ProcessResults(context XContext, writer gin.ResponseWriter, last bool, funcName string, results ...any)
}

// BusinessError business error
// Business errors are special cases in business scenarios that need to return different data types in different business contexts; essentially not considered errors, but an abstraction to facilitate business writing,
// allowing the same interface to have the ability to return different business codes and business data in special cases
type BusinessError = gone.BusinessError

type BindFieldFunc func(context *gin.Context, structVale reflect.Value) error
type BindStructFunc func(*gin.Context, reflect.Value) (reflect.Value, error)

type HttInjector interface {
	StartBindFuncs()
	BindFuncs() BindStructFunc
}

type Middleware interface {
	Process(ctx *gin.Context)
}

const (
	// IdGoneGin , IdGoneGinRouter , IdGoneGinProcessor, IdGoneGinProxy, IdGoneGinResponser, IdHttpInjector;
	// The GonerIds of Goners in goner/gin, which integrates gin framework for web request.
	IdGoneGin              = "gone-gin"
	IdGoneGinRouter        = "gone-gin-router"
	IdGoneGinSysMiddleware = "gone-gin-sys-middleware"
	IdGoneGinProxy         = "gone-gin-proxy"
	IdGoneGinResponser     = "gone-gin-responser"
	IdHttpInjector         = "http"
)
