package injector

import (
	"github.com/gone-io/gone/v2"
	"reflect"
)

type CompiledFunc[P any] func(P) ([]reflect.Value, error)

type BindFunc[P any] func(P) (reflect.Value, error)

type FieldSetter[P any] func(context P, fieldValue reflect.Value) error

type NameParser[P any] interface {

	// BuildParser 构建一个 BindFunc[P] 用于解析 结构体字段
	BuildParser(keyMap map[string]string, field reflect.StructField) (func(P) (reflect.Value, error), error)

	Name() string
}

type TypeParser[P any] interface {

	//Parse 从context中解析出 类型t的reflect.Value
	Parse(context P) (reflect.Value, error)
	Type() reflect.Type
}

type BindExecutor[P any] interface {
	InjectedByType(pt reflect.Type) BindFunc[P]
	FindFieldSetter(conf string, field reflect.StructField) (FieldSetter[P], error)
}

type DelayBindInjector[P any] interface {
	gone.StructFieldInjector
	Prepare(x any) (CompiledFunc[P], error)
}

var _ DelayBindInjector[any] = (*delayBindInjector[any])(nil)

type Field struct {
	s uintptr
}
