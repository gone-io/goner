package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"net/http"
	"reflect"
)

type httpRequestTypeParser struct {
	gone.Flag
}

func (c *httpRequestTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context.Request), nil
}

func (c *httpRequestTypeParser) Type() reflect.Type {
	return reflect.TypeOf((*http.Request)(nil))
}
