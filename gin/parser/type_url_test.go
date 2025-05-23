package parser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_urlTypeParser_Type(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Type
	}{
		{
			name: "test",
			want: reflect.TypeOf(&url.URL{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &urlTypeParser{}
			assert.Equalf(t, tt.want, c.Type(), "Type()")
		})
	}
}

func Test_urlTypeParser_Parse(t *testing.T) {

	address := &url.URL{
		Path: "/test",
	}

	type args struct {
		context *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test",
			args: args{
				context: &gin.Context{
					Request: &http.Request{
						URL: address,
					},
				},
			},
			want:    reflect.ValueOf(address),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &urlTypeParser{}
			got, err := c.Parse(tt.args.context)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.args.context)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Parse(%v)", tt.args.context)
		})
	}
}
