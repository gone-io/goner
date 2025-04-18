package g

import (
	"github.com/gone-io/gone/v2"
	"net"
	"reflect"
)

func Recover(logger gone.Logger) {
	if r := recover(); r != nil {
		logger.Errorf(
			"panic: %v, %s",
			r,
			gone.PanicTrace(2, 1),
		)
	}
}

func GetLocalIps() []net.IP {
	if addrs, err := net.InterfaceAddrs(); err != nil {
		panic(gone.ToErrorWithMsg(err, "cannot get ip address"))
	} else {
		var ips []net.IP
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
		return ips
	}
}

type LoadOp struct {
	goner   gone.Goner
	options []gone.Option
	f       gone.LoadFunc
}

func L(g gone.Goner, options ...gone.Option) *LoadOp {
	return &LoadOp{
		goner:   g,
		options: options,
	}
}

func F(loadFunc gone.LoadFunc) *LoadOp {
	return &LoadOp{
		f: loadFunc,
	}
}

func BuildOnceLoadFunc(ops ...*LoadOp) gone.LoadFunc {
	return gone.OnceLoad(func(loader gone.Loader) error {
		for _, op := range ops {
			if op.goner != nil {
				err := loader.Load(
					op.goner,
					op.options...,
				)
				if err != nil {
					return gone.ToError(err)
				}
			}
			if op.f != nil {
				err := op.f(loader)
				if err != nil {
					return gone.ToError(err)
				}
			}
		}
		return nil
	})
}

var m = make(map[any]struct{})

func SingLoadProviderFunc[P any, T any](fn gone.FunctionProvider[P, T], options ...gone.Option) gone.LoadFunc {
	return func(loader gone.Loader) error {
		if _, ok := m[&fn]; ok {
			return nil
		}
		m[&fn] = struct{}{}

		provider := gone.WrapFunctionProvider(fn)
		return loader.Load(provider, options...)
	}
}

func NamedThirdComponentLoadFunc[T any](name string, component T) gone.LoadFunc {
	return SingLoadProviderFunc(func(tagConf string, param struct{}) (T, error) {
		return component, nil
	}, gone.Name(name))
}

func GetComponentByName[T any](keeper gone.GonerKeeper, name string) (T, error) {
	component := keeper.GetGonerByName(name)
	if component == nil {
		return *new(T), gone.NewInnerError("not found", gone.GonerNameNotFound)
	}

	if t, ok := component.(T); ok {
		return t, nil
	}

	if g, ok := component.(gone.Provider[T]); ok {
		return g.Provide(name)
	}

	if g, ok := component.(gone.NoneParamProvider[T]); ok {
		return g.Provide()
	}

	if g, ok := component.(gone.NamedProvider); ok {
		provide, err := g.Provide(name, reflect.TypeOf(new(T)).Elem())
		if err != nil {
			return *new(T), gone.ToError(err)
		}
		return provide.(T), nil
	}
	return *new(T), gone.NewInnerError("not found compatible component", gone.GonerNameNotFound)
}
