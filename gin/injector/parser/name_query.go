package parser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"reflect"
)

type queryNameParser struct {
	gone.Flag
}

const anyName = "_*_"

var anyType = gone.GetInterfaceType(new(any))

func (s *queryNameParser) BuildParser(keyMap map[string]string, field reflect.StructField) (func(context *gin.Context) (reflect.Value, error), error) {
	t := field.Type
	mainKey := keyMap[s.Name()]

	if keyMap[anyName] == "true" {
		switch {
		case t.Kind() == reflect.String:
			return func(context *gin.Context) (reflect.Value, error) {
				query := context.Request.URL.Query()
				return reflect.ValueOf(query.Encode()), nil
			}, nil
		case t.Kind() == reflect.Struct || t.Kind() == reflect.Map || t.Kind() == reflect.Slice || t.Implements(anyType):
			return func(context *gin.Context) (reflect.Value, error) {
				value := reflect.New(t)
				if err := context.BindQuery(value.Interface()); err != nil {
					return emptyValue, gone.NewParameterError(fmt.Sprintf("bind query error: %s", err.Error()))
				}
				return value.Elem(), nil
			}, nil

		case t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct:
			return func(context *gin.Context) (reflect.Value, error) {
				value := reflect.New(t.Elem())
				if err := context.BindQuery(value.Interface()); err != nil {
					return emptyValue, gone.NewParameterError(fmt.Sprintf("bind query error: %s", err.Error()))
				}
				return value, nil
			}, nil

		default:
			return nil, gone.NewInnerError(fmt.Sprintf("build parser failed for field(name=%s)", field.Name), gone.InjectError)
		}
	} else {
		parser, err := BuildParser(t)
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("build parser failed for field(name=%s)", field.Name))
		}

		return func(context *gin.Context) (reflect.Value, error) {
			param := context.Query(mainKey)

			if v, err := parser(param); err != nil {
				return emptyValue, gone.NewParameterError(fmt.Sprintf("parse cooke[name=%s] error: %s", mainKey, err.Error()))
			} else {
				return v, nil
			}
		}, nil
	}
}

func (s *queryNameParser) Name() string {
	return "query"
}
