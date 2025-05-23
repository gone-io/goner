package parser

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
)

func Test_headerNameParser_BuildParser(t *testing.T) {
	type Q struct {
		F1 string `form:"f1"`
		X  int    `form:"x"`
	}

	type X struct {
		F1  string   `gone:"http,header=*"`
		F2  Q        `gone:"http,header"`
		F3  *Q       `gone:"http,header"`
		F4  chan int `gone:"http,header"`
		F5  string   `gone:"http,header=key"`
		f6  float64  `gone:"http,header"`
		f7  string   `gone:"http,header"`
		f8  []string `gone:"http,header=key"`
		f9  []int    `gone:"http,header"`
		f10 []func() `gone:"http,header"`
	}

	xType := reflect.TypeOf(&X{}).Elem()

	tests := []struct {
		name     string
		field    string
		buildErr assert.ErrorAssertionFunc
		header   http.Header
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
			header: http.Header{
				"Key": []string{"1"},
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse one header as number",
			field:    "f6",
			buildErr: assert.NoError,
			header: http.Header{
				"F6": []string{"1.1"},
			},
			wantErr: assert.NoError,
			want:    1.1,
		},
		{
			name:     "parse one header as number and parse error",
			field:    "f6",
			buildErr: assert.NoError,
			header: http.Header{
				"F6": []string{"1..1"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:     "parse one header as string use filed name as key",
			field:    "f7",
			buildErr: assert.NoError,
			header: http.Header{
				"F7": []string{"1"},
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse header not found",
			field:    "f6",
			buildErr: assert.NoError,
			header: http.Header{
				"F6": []string{"p1"},
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

			s := &headerNameParser{}
			parser, err := s.BuildParser(keyMap, field)
			if !tt.buildErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.NotNil(t, parser)

			value, err := parser(&gin.Context{
				Request: &http.Request{
					Header: tt.header,
				},
			})
			if !tt.wantErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.True(t, reflect.DeepEqual(value.Interface(), tt.want))
		})
	}
}

func Test_headerNameParser_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test",
			want: "header",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := headerNameParser{}
			assert.Equalf(t, tt.want, b.Name(), "Name()")
		})
	}
}
