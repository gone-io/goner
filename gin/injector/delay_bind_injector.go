package injector

import (
	"github.com/gone-io/gone/v2"
	"reflect"
)

type delayBindInjector[P any] struct {
	gone.Flag
	fieldBindMap map[Field]FieldSetter[P]
	funcInjector gone.FuncInjector `gone:"*"`
	bindExecutor BindExecutor[P]   `gone:"*"`

	name string
}

func (h *delayBindInjector[P]) GonerName() string {
	return h.name
}

func (h *delayBindInjector[P]) startBindFuncs() {
	h.fieldBindMap = map[Field]FieldSetter[P]{}
}

func (h *delayBindInjector[P]) setFieldBindMap(fieldValue reflect.Value, fn FieldSetter[P]) {
	h.fieldBindMap[Field{
		s: fieldValue.UnsafeAddr(),
	}] = fn
}

func (h *delayBindInjector[P]) getFieldBindMap(fieldValue reflect.Value) FieldSetter[P] {
	return h.fieldBindMap[Field{
		s: fieldValue.UnsafeAddr(),
	}]
}

func (h *delayBindInjector[P]) Inject(tagConf string, field reflect.StructField, fieldValue reflect.Value) error {
	if fn, err := h.bindExecutor.FindFieldSetter(tagConf, field); err != nil {
		return err
	} else {
		h.setFieldBindMap(fieldValue, fn)
	}
	return nil
}

func (h *delayBindInjector[P]) Prepare(f any) (CompiledFunc[P], error) {
	h.startBindFuncs()

	values, err := h.funcInjector.InjectFuncParameters(
		f,
		func(pt reflect.Type, i int, injected bool) any {
			bindFn := h.bindExecutor.InjectedByType(pt)
			if bindFn == nil || reflect.ValueOf(bindFn).IsNil() {
				return nil
			}
			return bindFn
		},
		nil,
	)
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "inject func failed")
	}

	of := reflect.ValueOf(f)
	link := h.link(values)

	return func(context P) ([]reflect.Value, error) {
		if args, err := link(context); err != nil {
			return nil, err
		} else {
			return of.Call(args), nil
		}
	}, nil
}

func (h *delayBindInjector[P]) link(values []reflect.Value) func(context P) (list []reflect.Value, err error) {
	var binds []BindFunc[P]
	var BindFuncType = reflect.TypeOf((*BindFunc[P])(nil)).Elem()

	for _, v := range values {
		if v.Type() == BindFuncType {
			a := v.Interface()
			bind := a.(BindFunc[P])
			binds = append(binds, bind)
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			bindFunc := h.buildStructBindFunc(v)
			binds = append(binds, bindFunc)
		case reflect.Ptr:
			pointValue := v.Elem()

			if pointValue.Type().Kind() == reflect.Struct {
				bindFunc := h.buildStructBindFunc(pointValue)
				binds = append(binds, bindFunc)
				continue
			}
			fallthrough
		default:
			binds = append(binds, func(context P) (reflect.Value, error) {
				return v, nil
			})
		}
	}

	return func(context P) (list []reflect.Value, err error) {
		for _, f := range binds {
			var value reflect.Value
			if value, err = f(context); err != nil {
				return nil, err
			} else {
				list = append(list, value)
			}
		}
		return
	}
}

func (h *delayBindInjector[P]) buildStructBindFunc(structValue reflect.Value) BindFunc[P] {
	structType := structValue.Type()
	fieldNum := structType.NumField()
	m := make(map[int]FieldSetter[P])

	for i := 0; i < fieldNum; i++ {
		value := structValue.Field(i)
		setter := h.getFieldBindMap(value)
		if setter != nil {
			m[i] = setter
		}
	}
	return func(context P) (reflect.Value, error) {
		newValue := reflect.New(structType).Elem()
		newValue.Set(structValue)
		for i, setter := range m {
			fieldValue := newValue.Field(i)
			if !fieldValue.CanSet() {
				fieldValue = gone.BlackMagic(fieldValue)
			}
			if err := setter(context, fieldValue); err != nil {
				return reflect.Value{}, err
			}
		}
		return newValue, nil
	}
}
