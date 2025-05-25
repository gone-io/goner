package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_originContextTypeParser_Type(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Type
	}{
		{
			name: "test",
			want: reflect.TypeOf((*gin.Context)(nil)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ginContextTypeParser{}
			assert.Equalf(t, tt.want, c.Type(), "Type()")
		})
	}
}

func Test_originContextTypeParser_Parse(t *testing.T) {
	ctx := &gin.Context{}

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
				context: ctx,
			},
			want:    reflect.ValueOf(ctx),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ginContextTypeParser{}
			got, err := c.Parse(tt.args.context)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.args.context)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Parse(%v)", tt.args.context)
		})
	}
}
