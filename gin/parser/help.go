package parser

import (
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone/v2"
	"reflect"
)

var emptyValue = reflect.Value{}

func BuildParser(t reflect.Type) (func(str string) (reflect.Value, error), error) {
	switch t.Kind() {
	case reflect.String:
		return func(str string) (reflect.Value, error) {
			return reflect.ValueOf(str), nil
		}, nil
	case reflect.Struct, reflect.Slice, reflect.Map:
		return func(str string) (reflect.Value, error) {
			value := reflect.New(t)
			err := json.Unmarshal([]byte(str), value.Interface())
			if err != nil {
				return emptyValue, gone.NewParameterError(fmt.Sprintf("cannot parse json: %s; %s", str, err.Error()))
			}
			return value.Elem(), nil
		}, nil
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool:
		return func(str string) (reflect.Value, error) {
			value := reflect.New(t)
			err := gone.SetValueByReflect(value, str)
			if err != nil {
				return emptyValue, gone.NewParameterError(fmt.Sprintf("cannot parse value: %s; %s", str, err.Error()))
			}
			return value.Elem(), nil
		}, nil

	default:
		return nil, gone.NewInnerError(fmt.Sprintf("unsupported type: %s", gone.GetTypeName(t)), gone.InjectError)
	}
}
