package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"net/http"
	"reflect"
)

// http.ResponseWriter
type httpResponseTypeParser struct {
	gone.Flag
}

func (c *httpResponseTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context.Writer), nil
}

func (c *httpResponseTypeParser) Type() reflect.Type {
	return gone.GetInterfaceType(new(http.ResponseWriter))
}
