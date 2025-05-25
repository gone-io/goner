package parser

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_httpHeaderTypeParser_Type(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Type
	}{
		{
			name: "test",
			want: reflect.TypeOf(http.Header{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpHeaderTypeParser{}
			assert.Equalf(t, tt.want, c.Type(), "Type()")
		})
	}
}

func Test_httpHeaderTypeParser_Parse(t *testing.T) {
	header := http.Header{
		"Content-Type": []string{"application/json"},
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
						Header: header,
					},
				},
			},
			want:    reflect.ValueOf(header),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpHeaderTypeParser{}
			got, err := c.Parse(tt.args.context)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.args.context)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Parse(%v)", tt.args.context)
		})
	}
}
