package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

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

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(nameParser []NameParser[*gin.Context], typeParsers []TypeParser[*gin.Context]) {
			assert.Equal(t, 5, len(nameParser))
			assert.Equal(t, 5, len(typeParsers))
		})
}
