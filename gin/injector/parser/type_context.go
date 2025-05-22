package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	goneGin "github.com/gone-io/goner/gin"
	"reflect"
)

// *goneGin.Context
type contextTypeParser struct {
	gone.Flag
}

func (c *contextTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	goneCtx := goneGin.Context{
		Context: context,
	}

	return reflect.ValueOf(&goneCtx), nil
}

func (c *contextTypeParser) Type() reflect.Type {
	return reflect.TypeOf((*goneGin.Context)(nil))
}
