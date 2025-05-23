package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
)

func Test_paramNameParser_BuildParser(t *testing.T) {
	type Q struct {
		F1 string `form:"f1"`
		X  int    `form:"x"`
	}

	type X struct {
		F1  string   `gone:"http,param=*"`
		F2  Q        `gone:"http,param"`
		F3  *Q       `gone:"http,param"`
		F4  chan int `gone:"http,param"`
		F5  string   `gone:"http,param=key"`
		f6  float64  `gone:"http,param"`
		f7  string   `gone:"http,param"`
		f8  []string `gone:"http,param=key"`
		f9  []int    `gone:"http,param"`
		f10 []func() `gone:"http,param"`
	}

	xType := reflect.TypeOf(&X{}).Elem()

	tests := []struct {
		name     string
		field    string
		buildErr assert.ErrorAssertionFunc
		params   gin.Params
		wantErr  assert.ErrorAssertionFunc
		want     any
	}{
		{
			name:  "parse one query as unsupported slice",
			field: "F4",
			buildErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:     "parse string",
			field:    "F5",
			buildErr: assert.NoError,
			params: gin.Params{
				{Key: "key", Value: "1"},
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse one param as number",
			field:    "f6",
			buildErr: assert.NoError,
			params: gin.Params{
				{Key: "f6", Value: "1.1"},
			},
			wantErr: assert.NoError,
			want:    1.1,
		},
		{
			name:     "parse one param as number and parse error",
			field:    "f6",
			buildErr: assert.NoError,
			params: gin.Params{
				{Key: "f6", Value: "1..1"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:     "parse one param as string use filed name as key",
			field:    "f7",
			buildErr: assert.NoError,
			params: gin.Params{
				{Key: "f7", Value: "1"},
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse param not found",
			field:    "f6",
			buildErr: assert.NoError,
			params: gin.Params{
				{Key: "f6", Value: "p1"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := xType.FieldByName(tt.field)
			assert.True(t, ok)

			goneTag := field.Tag.Get("gone")
			_, extend := gone.ParseGoneTag(goneTag)
			keyMap, keys := gone.TagStringParse(extend)
			name := keys[0]
			if keyMap[name] == "" {
				keyMap[anyName] = "true"
				keyMap[name] = field.Name
			}

			s := &paramNameParser{}
			parser, err := s.BuildParser(keyMap, field)
			if !tt.buildErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.NotNil(t, parser)

			value, err := parser(&gin.Context{
				Params: tt.params,
			})
			if !tt.wantErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.True(t, reflect.DeepEqual(value.Interface(), tt.want))
		})
	}
}

func Test_paramNameParser_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test",
			want: "param",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := paramNameParser{}
			assert.Equalf(t, tt.want, b.Name(), "Name()")
		})
	}
}
