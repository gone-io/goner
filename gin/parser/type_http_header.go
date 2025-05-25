package parser

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"net/http"
	"reflect"
)

type httpHeaderTypeParser struct {
	gone.Flag
}

func (c *httpHeaderTypeParser) Parse(context *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(context.Request.Header), nil
}

func (c *httpHeaderTypeParser) Type() reflect.Type {
	return reflect.TypeOf(http.Header{})
}
