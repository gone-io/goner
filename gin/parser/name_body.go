package parser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"io"
	"net/http"
	"reflect"
)

// for body parser
type bodyNameParser struct {
	gone.Flag
}

var bytesType = reflect.TypeOf([]byte{})
var readerType = gone.GetInterfaceType(new(io.Reader))
var readCloserType = gone.GetInterfaceType(new(io.ReadCloser))

func (b bodyNameParser) BuildParser(_ map[string]string, field reflect.StructField) (func(context *gin.Context) (reflect.Value, error), error) {
	t := field.Type
	switch {
	case t == bytesType:
		return func(context *gin.Context) (reflect.Value, error) {
			all, err := io.ReadAll(context.Request.Body)
			if err != nil {
				return emptyValue, gone.NewParameterError(fmt.Sprintf("bind body error: %s", err.Error()))
			}
			return reflect.ValueOf(all), nil
		}, nil

	case t.Implements(readerType) || t.Implements(readCloserType):
		return func(context *gin.Context) (reflect.Value, error) {
			return reflect.ValueOf(context.Request.Body), nil
		}, nil

	case t == anyType || t.Kind() == reflect.Struct || t.Kind() == reflect.Map || t.Kind() == reflect.Slice:
		return func(context *gin.Context) (reflect.Value, error) {
			value := reflect.New(t)
			if err := context.ShouldBind(value.Interface()); err != nil {
				return emptyValue, gone.NewParameterError(fmt.Sprintf("bind body error: %s", err.Error()))
			}
			return value.Elem(), nil
		}, nil

	default:
		switch t.Kind() {
		case reflect.String:
			return func(context *gin.Context) (reflect.Value, error) {
				all, err := io.ReadAll(context.Request.Body)
				if err != nil {
					return emptyValue, gone.NewParameterError(fmt.Sprintf("bind body error: %s", err.Error()))
				}
				return reflect.ValueOf(string(all)), nil
			}, nil
		case reflect.Ptr:
			if t.Elem().Kind() == reflect.Struct {
				return func(context *gin.Context) (reflect.Value, error) {
					value := reflect.New(t.Elem())
					if err := context.ShouldBind(value.Interface()); err != nil {
						return emptyValue, gone.NewParameterError(fmt.Sprintf("bind body error: %s", err.Error()))
					}
					return value, nil
				}, nil
			}
			fallthrough
		default:
			return nil, gone.NewInnerError(fmt.Sprintf("unsupported type %s(field=%s) ", gone.GetTypeName(t), field.Name), http.StatusInternalServerError)
		}
	}
}

func (b bodyNameParser) Name() string {
	return "body"
}
