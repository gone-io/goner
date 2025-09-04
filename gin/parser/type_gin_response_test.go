package parser

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_responseTypeParser_Type(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Type
	}{
		{
			name: "test",
			want: gone.GetInterfaceType(new(gin.ResponseWriter)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &responseTypeParser{}
			assert.Equalf(t, tt.want, c.Type(), "Type()")
		})
	}
}

func Test_responseTypeParser_Parse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

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
			want:    reflect.ValueOf(ctx.Writer),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &responseTypeParser{}
			got, err := c.Parse(tt.args.context)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.args.context)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Parse(%v)", tt.args.context)
		})
	}
}
