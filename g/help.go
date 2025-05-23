package g

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"net"
	"net/http"
	"reflect"
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
	addrs, err := net.InterfaceAddrs()
	PanicIfErr(err)
	var ips []net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP)
		}
	}
	return ips
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
	return func(loader gone.Loader) error {
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
	}
}

// SingLoadProviderFunc creates a loading function for singleton Provider
// P: Provider's parameter type
// T: Component type provided by Provider
// fn: Function to create components
// options: Loading options
// Returns: Loading function that ensures single loading
// Deprecated since v2.1.0, use gone.BuildSingProviderLoadFunc instead:
func SingLoadProviderFunc[P any, T any](fn gone.FunctionProvider[P, T], options ...gone.Option) gone.LoadFunc {
	return gone.BuildSingProviderLoadFunc(fn, options...)
}

// NamedThirdComponentLoadFunc creates a named third-party component loading function
// T: Component type
// name: Component name
// component: Third-party component instance
// Returns: Loading function for third-party components
func NamedThirdComponentLoadFunc[T any](name string, component T) gone.LoadFunc {
	return gone.BuildThirdComponentLoadFunc(component, gone.Name(name))
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

var appMap = make(map[string]*gone.Application)

// App creates or retrieves an application instance with the specified name and loading functions
// name: Application name
// loadFuncs: Loading functions
// Returns: Application instance
func App(name string, loadFuncs ...gone.LoadFunc) (app *gone.Application) {
	var ok bool
	if app, ok = appMap[name]; !ok {
		app = gone.NewApp(loadFuncs...)
		appMap[name] = app
	}
	return app
}

// PanicIfErr panics if an error occurs
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ResultError return wrapped error if error occurs, otherwise returns the result
func ResultError[T any](t *T, err error, msg string) (*T, error) {
	if err != nil {
		return nil, gone.NewInnerErrorSkip(fmt.Sprintf("%s: %v", msg, err), http.StatusInternalServerError, 3)
	}
	return t, nil
}

func ErrorPrinter(logger gone.Logger, err error, msg string, args ...any) {
	if err != nil {
		logger.Errorf("%v", gone.ToErrorWithMsg(err, fmt.Sprintf(msg, args...)))
	}
}

var m = make(map[gone.Loader]map[gone.Goner]struct{})

func IsLoaded(loader gone.Loader, goner gone.Goner) bool {
	if m2, ok := m[loader]; ok {
		if _, ok := m2[goner]; ok {
			return true
		}
	} else {
		m[loader] = make(map[gone.Goner]struct{})
	}
	m[loader][goner] = struct{}{}
	return false
}
