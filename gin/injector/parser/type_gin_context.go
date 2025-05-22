package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"reflect"
)

// *gin.Context
type originContextTypeParser struct {
	gone.Flag
}

func (c *originContextTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context), nil
}

func (c *originContextTypeParser) Type() reflect.Type {
	return reflect.TypeOf((*gin.Context)(nil))
}
