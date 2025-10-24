package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/gin/injector"
	"reflect"
	"time"
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
	prepare, err := p.injector.Prepare(x)
	g.PanicIfErr(err)

	return func(context *gin.Context) {
		if p.stat {
			defer TimeStat(funcName+"-inject-proxy", time.Now(), p.log.Infof)
		}
		values, err := prepare(context)
		if err != nil {
			p.responser.Failed(context, err)
			context.Abort()
			return
		}
		p.resultProcess(values, context, funcName, last)
	}
}

func (p *proxy) resultProcess(values []reflect.Value, context *gin.Context, funcName string, last bool) {
	var results []any
	for i := 0; i < len(values); i++ {
		arg := values[i]
		switch arg.Kind() {
		case reflect.Invalid:
			//results = append(results, nil)
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
			if arg.IsNil() {
				//results = append(results, nil)
				continue
			}
			fallthrough
		default:
			results = append(results, arg.Interface())
		}
	}
	p.responser.ProcessResults(context, context.Writer, last, funcName, results...)
}
