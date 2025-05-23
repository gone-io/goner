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

func Test_cookeNameParser_BuildParser(t *testing.T) {
	type Q struct {
		F1 string `form:"f1"`
		X  int    `form:"x"`
	}

	type X struct {
		F1  string   `gone:"http,cookie=*"`
		F2  Q        `gone:"http,cookie"`
		F3  *Q       `gone:"http,cookie"`
		F4  chan int `gone:"http,cookie"`
		F5  string   `gone:"http,cookie=key"`
		f6  float64  `gone:"http,cookie"`
		f7  string   `gone:"http,cookie"`
		f8  []string `gone:"http,cookie=key"`
		f9  []int    `gone:"http,cookie"`
		f10 []func() `gone:"http,cookie"`
	}

	xType := reflect.TypeOf(&X{}).Elem()

	tests := []struct {
		name     string
		field    string
		buildErr assert.ErrorAssertionFunc
		cookie   *http.Cookie
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
			cookie: &http.Cookie{
				Name:  "key",
				Value: "1",
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse one cookie as number",
			field:    "f6",
			buildErr: assert.NoError,
			cookie: &http.Cookie{
				Name:  "f6",
				Value: "1.1",
			},
			wantErr: assert.NoError,
			want:    1.1,
		},
		{
			name:     "parse one cookie as number and parse error",
			field:    "f6",
			buildErr: assert.NoError,
			cookie: &http.Cookie{
				Name:  "f6",
				Value: "1..1",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:     "parse one cookie as string use filed name as key",
			field:    "f7",
			buildErr: assert.NoError,
			cookie: &http.Cookie{
				Name:  "f7",
				Value: "1",
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse cookie not found",
			field:    "f7",
			buildErr: assert.NoError,
			cookie:   nil,
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

			s := &cookeNameParser{}
			parser, err := s.BuildParser(keyMap, field)
			if !tt.buildErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.NotNil(t, parser)

			req := &http.Request{
				Header: http.Header{},
			}
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			value, err := parser(&gin.Context{
				Request: req,
			})
			if !tt.wantErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.True(t, reflect.DeepEqual(value.Interface(), tt.want))
		})
	}
}

func Test_cookeNameParser_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test",
			want: "cookie",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := cookeNameParser{}
			assert.Equalf(t, tt.want, b.Name(), "Name()")
		})
	}
}
