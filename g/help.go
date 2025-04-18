package g

import (
	"net"
	"reflect"

	"github.com/gone-io/gone/v2"
)

// Recover captures and logs panics to prevent program crashes
// logger: Logger for recording panic information
func Recover(logger gone.Logger) {
	if r := recover(); r != nil {
		logger.Errorf(
			"panic: %v, %s",
			r,
			gone.PanicTrace(2, 1),
		)
	}
}

// GetLocalIps gets all non-loopback IPv4 addresses of the local machine
// Returns: List of all available IPv4 addresses
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

// LoadOp struct encapsulates loading operations
// goner: Component to be loaded
// options: Loading options
// f: Loading function
type LoadOp struct {
	goner   gone.Goner
	options []gone.Option
	f       gone.LoadFunc
}

// L creates a LoadOp for loading components
// g: Component to be loaded
// options: Loading options
// Returns: New LoadOp instance
func L(g gone.Goner, options ...gone.Option) *LoadOp {
	return &LoadOp{
		goner:   g,
		options: options,
	}
}

// F creates a LoadOp containing a loading function
// loadFunc: Custom loading function
// Returns: New LoadOp instance
func F(loadFunc gone.LoadFunc) *LoadOp {
	return &LoadOp{
		f: loadFunc,
	}
}

// BuildOnceLoadFunc builds a loading function that executes only once
// ops: List of LoadOps to execute
// Returns: Loading function that ensures single execution
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

// SingLoadProviderFunc creates a loading function for singleton Provider
// P: Provider's parameter type
// T: Component type provided by Provider
// fn: Function to create components
// options: Loading options
// Returns: Loading function that ensures single loading
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

// NamedThirdComponentLoadFunc creates a named third-party component loading function
// T: Component type
// name: Component name
// component: Third-party component instance
// Returns: Loading function for third-party components
func NamedThirdComponentLoadFunc[T any](name string, component T) gone.LoadFunc {
	return SingLoadProviderFunc(func(tagConf string, param struct{}) (T, error) {
		return component, nil
	}, gone.Name(name))
}

// GetComponentByName gets a component of specified type by name
// T: Component type to get
// keeper: Component manager
// name: Component name
// Returns: Found component instance and possible error
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
