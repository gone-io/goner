package parser

import (
	"fmt"
	goneGin "github.com/gone-io/goner/gin"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_contextTypeParser_Type(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Type
	}{
		{
			name: "test",
			want: reflect.TypeOf((*goneGin.Context)(nil)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &contextTypeParser{}
			assert.Equalf(t, tt.want, c.Type(), "Type()")
		})
	}
}

func Test_contextTypeParser_Parse(t *testing.T) {
	ctx := &gin.Context{}

	type args struct {
		context *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test",
			args: args{
				context: ctx,
			},
			want:    &goneGin.Context{Context: ctx},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &contextTypeParser{}
			got, err := c.Parse(tt.args.context)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.args.context)) {
				return
			}
			assert.True(t, reflect.DeepEqual(tt.want, got.Interface()))
		})
	}
}
