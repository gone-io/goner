package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"reflect"
)

// gin.ResponseWriter
type responseTypeParser struct {
	gone.Flag
}

func (c *responseTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context.Writer), nil
}

func (c *responseTypeParser) Type() reflect.Type {
	return reflect.TypeOf((gin.ResponseWriter)(nil))
}
