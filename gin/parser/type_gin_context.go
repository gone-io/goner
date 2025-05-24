package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"reflect"
)

// *gin.Context
type ginContextTypeParser struct {
	gone.Flag
}

func (c *ginContextTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context), nil
}

func (c *ginContextTypeParser) Type() reflect.Type {
	return reflect.TypeOf((*gin.Context)(nil))
}
