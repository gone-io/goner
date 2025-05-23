package parser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_httpRequestTypeParser_Type(t *testing.T) {

	tests := []struct {
		name string
		want reflect.Type
	}{
		{
			name: "test",
			want: reflect.TypeOf(&http.Request{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpRequestTypeParser{}
			assert.Equalf(t, tt.want, c.Type(), "Type()")
		})
	}
}

func Test_httpRequestTypeParser_Parse(t *testing.T) {

	request := http.Request{}

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
					Request: &request,
				},
			},
			want:    reflect.ValueOf(&request),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpRequestTypeParser{}
			got, err := c.Parse(tt.args.context)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.args.context)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Parse(%v)", tt.args.context)
		})
	}
}
