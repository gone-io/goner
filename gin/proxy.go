package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/gin/injector"
	"reflect"
)

type proxy struct {
	gone.Flag
	log          gone.Logger                              `gone:"*"`
	funcInjector gone.FuncInjector                        `gone:"*"`
	responser    Responser                                `gone:"*"`
	injector     injector.DelayBindInjector[*gin.Context] `gone:"*"`
	stat         bool                                     `gone:"config,server.proxy.stat,default=false"`
}

func (p *proxy) GonerName() string {
	return IdGoneGinProxy
}

func (p *proxy) Proxy(handlers ...HandlerFunc) (arr []gin.HandlerFunc) {
	count := len(handlers)
	for i := 0; i < count-1; i++ {
		arr = append(arr, p.proxyOne(handlers[i], false))
	}
	arr = append(arr, p.proxyOne(handlers[count-1], true))
	return arr
}

func (p *proxy) ProxyForMiddleware(handlers ...HandlerFunc) (arr []gin.HandlerFunc) {
	count := len(handlers)
	for i := 0; i < count; i++ {
		arr = append(arr, p.proxyOne(handlers[i], false))
	}
	return arr
}

func (p *proxy) proxyOne(x HandlerFunc, last bool) gin.HandlerFunc {
	funcName := gone.GetFuncName(x)

	switch f := x.(type) {
	case func(*Context) (any, error):
		return func(context *gin.Context) {
			data, err := f(&Context{Context: context})
			p.responser.ProcessResults(context, context.Writer, last, funcName, data, err)
		}
	case func(*Context) error:
		return func(context *gin.Context) {
			err := f(&Context{Context: context})
			p.responser.ProcessResults(context, context.Writer, last, funcName, err)
		}
	case func(*Context):
		return func(context *gin.Context) {
			f(&Context{Context: context})
			p.responser.ProcessResults(context, context.Writer, last, funcName)
		}
	case func(ctx *gin.Context):
		return x.(func(ctx *gin.Context))

	case func(ctx *gin.Context) (any, error):
		return func(context *gin.Context) {
			data, err := f(context)
			p.responser.ProcessResults(context, context.Writer, last, funcName, data, err)
		}
	case func(ctx *gin.Context) error:
		return func(context *gin.Context) {
			err := f(context)
			p.responser.ProcessResults(context, context.Writer, last, funcName, err)
		}
	case func():
		return func(context *gin.Context) {
			f()
			p.responser.ProcessResults(context, context.Writer, last, funcName)
		}
	case func() (any, error):
		return func(context *gin.Context) {
			data, err := f()
			p.responser.ProcessResults(context, context.Writer, last, funcName, data, err)
		}
	case func() error:
		return func(context *gin.Context) {
			err := f()
			p.responser.ProcessResults(context, context.Writer, last, funcName, err)
		}
	default:
		return p.buildProxyFn(x, funcName, last)
	}
}

func (p *proxy) resultProcess(values []reflect.Value, context *gin.Context, funcName string, last bool) {
	var results []any
	for i := 0; i < len(values); i++ {
		arg := values[i]

		if arg.Kind() == reflect.Interface {
			elem := arg.Elem()
			switch elem.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
				if elem.IsNil() {
					results = append(results, nil)
					continue
				}
			default:
			}
		}
		results = append(results, arg.Interface())
	}
	p.responser.ProcessResults(context, context.Writer, last, funcName, results...)
}

func (p *proxy) buildProxyFn(x HandlerFunc, funcName string, last bool) gin.HandlerFunc {
	prepare, err := p.injector.Prepare(x)
	g.PanicIfErr(err)

	return func(context *gin.Context) {
		values, err := prepare(context)
		if err != nil {
			p.responser.Failed(context, err)
		}

		p.resultProcess(values, context, funcName, last)
	}
}
