package parser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_queryNameParser_BuildParser(t *testing.T) {
	type Q struct {
		F1 string `form:"f1"`
		X  int    `form:"x"`
	}

	type X struct {
		F1  string   `gone:"http,query=*"`
		F2  Q        `gone:"http,query"`
		F3  *Q       `gone:"http,query"`
		F4  chan int `gone:"http,query"`
		F5  string   `gone:"http,query=key"`
		f6  float64  `gone:"http,query"`
		f7  string   `gone:"http,query"`
		f8  []string `gone:"http,query=key"`
		f9  []int    `gone:"http,query"`
		f10 []func() `gone:"http,query"`
	}

	xType := reflect.TypeOf(&X{}).Elem()

	tests := []struct {
		name     string
		field    string
		buildErr assert.ErrorAssertionFunc
		url      url.URL
		wantErr  assert.ErrorAssertionFunc
		want     any
	}{
		{
			name:     "parse all query as string",
			field:    "F1",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f1=1&x=200",
			},
			wantErr: assert.NoError,
			want:    "f1=1&x=200",
		},
		{
			name:     "parse all query as struct",
			field:    "F2",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f1=1&x=200",
			},
			wantErr: assert.NoError,
			want: Q{
				F1: "1",
				X:  200,
			},
		},
		{
			name:     "parse all query as struct and parse err",
			field:    "F2",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f1=1&x=2l00",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:     "parse all query as pointer struct",
			field:    "F3",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f1=1&x=200",
			},
			wantErr: assert.NoError,
			want: &Q{
				F1: "1",
				X:  200,
			},
		},
		{
			name:     "parse all query as pointer struct and error",
			field:    "F3",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f1=1&x=2o00",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:  "parse all query and unsupported type",
			field: "F4",
			buildErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:     "parse one query as string",
			field:    "F5",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "key=1",
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse one query as number",
			field:    "f6",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f6=1.1",
			},
			wantErr: assert.NoError,
			want:    1.1,
		},
		{
			name:     "parse one query as number and parse error",
			field:    "f6",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f6=1..1",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:     "parse one query as string use filed name as key",
			field:    "f7",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f7=1",
			},
			wantErr: assert.NoError,
			want:    "1",
		},
		{
			name:     "parse one query as string slice",
			field:    "f8",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "key=1&key=2",
			},
			wantErr: assert.NoError,
			want:    []string{"1", "2"},
		},
		{
			name:     "parse one query as number slice",
			field:    "f9",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f9=1&f9=2",
			},
			wantErr: assert.NoError,
			want:    []int{1, 2},
		},
		{
			name:     "parse one query as number slice and parse err",
			field:    "f9",
			buildErr: assert.NoError,
			url: url.URL{
				RawQuery: "f9=1&f9=o12",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name:  "parse one query as unsupported slice",
			field: "f10",
			buildErr: func(t assert.TestingT, err error, i ...interface{}) bool {
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

			s := &queryNameParser{}
			parser, err := s.BuildParser(keyMap, field)
			if !tt.buildErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.NotNil(t, parser)
			value, err := parser(&gin.Context{
				Request: &http.Request{
					URL: &tt.url,
				},
			})
			if !tt.wantErr(t, err, fmt.Sprintf("BuildParser(%v, %v)", keyMap, tt.field)) {
				return
			}
			assert.True(t, reflect.DeepEqual(value.Interface(), tt.want))
		})
	}
}
