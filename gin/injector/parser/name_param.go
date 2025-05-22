package parser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"reflect"
)

type paramNameParser struct {
	gone.Flag
}

func (s *paramNameParser) BuildParser(keyMap map[string]string, field reflect.StructField) (func(context *gin.Context) (reflect.Value, error), error) {
	t := field.Type
	mainKey := keyMap[s.Name()]

	parser, err := BuildParser(t)
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("build parser failed for field(name=%s)", field.Name))
	}
	return func(context *gin.Context) (reflect.Value, error) {
		param := context.Param(mainKey)

		if v, err := parser(param); err != nil {
			return emptyValue, gone.NewParameterError(fmt.Sprintf("parse cooke[name=%s] error: %s", mainKey, err.Error()))
		} else {
			return v, nil
		}
	}, nil
}

func (s *paramNameParser) Name() string {
	return "param"
}
