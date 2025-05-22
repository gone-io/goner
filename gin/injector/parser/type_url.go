package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"net/url"
	"reflect"
)

type urlTypeParser struct {
	gone.Flag
}

func (c *urlTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context.Request.URL), nil
}
func (c *urlTypeParser) Type() reflect.Type {
	return reflect.TypeOf((*url.URL)(nil)).Elem()
}
