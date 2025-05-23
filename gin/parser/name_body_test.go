package parser

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type errReader struct {
}

func (e *errReader) Close() error {
	return errors.New("test")
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test")
}

func Test_bodyNameParser_BuildParser(t *testing.T) {
	type Req struct {
		X int    `json:"x,omitempty"`
		Y string `json:"y,omitempty"`
	}

	type IN struct {
		req1 Req            `gone:"http,body"`
		req2 io.Reader      `gone:"http,body"`
		req3 []byte         `gone:"http,body"`
		req4 *Req           `gone:"http,body"`
		req5 string         `gone:"http,body"`
		req6 any            `gone:"http,body"`
		req7 []any          `gone:"http,body"`
		req8 map[string]any `gone:"http,body"`
		req9 *string        `gone:"http,body"`
	}

	var in IN

	of := reflect.TypeOf(&in).Elem()
	req1, _ := of.FieldByName("req1")
	req2, _ := of.FieldByName("req2")
	req3, _ := of.FieldByName("req3")
	req4, _ := of.FieldByName("req4")
	req5, _ := of.FieldByName("req5")
	req6, _ := of.FieldByName("req6")
	req7, _ := of.FieldByName("req7")
	req8, _ := of.FieldByName("req8")
	req9, _ := of.FieldByName("req9")

	type args struct {
		field   reflect.StructField
		context *gin.Context
	}
	tests := []struct {
		name     string
		args     args
		want     any
		buildErr bool
		parseErr bool
	}{
		{
			name: "parse struct",
			args: args{
				field: req1,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: Req{
				X: 1,
				Y: "2",
			},
		},
		{
			name: "parse io.reader",
			args: args{
				field: req2,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: func(r any) bool {
				if reader, ok := r.(io.Reader); !ok {
					return false
				} else {
					var buf = make([]byte, 15)

					n, err := reader.Read(buf)
					assert.Nil(t, err)
					assert.Equal(t, n, 15)
					assert.Equal(t, string(buf), `{"x":1,"y":"2"}`)
					return true
				}
			},
		},
		{
			name: "parse []byte",
			args: args{
				field: req3,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: []byte(`{"x":1,"y":"2"}`),
		},
		{
			name: "parse []byte failed",
			args: args{
				field: req3,
				context: &gin.Context{
					Request: &http.Request{
						Body: &errReader{},
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			parseErr: true,
		},
		{
			name: "parse struct ptr",
			args: args{
				field: req4,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: func(r any) bool {
				if req, ok := r.(*Req); !ok {
					return false
				} else {
					return req.X == 1 && req.Y == "2"
				}
			},
		},
		{
			name: "parse struct ptr failed",
			args: args{
				field: req4,
				context: &gin.Context{
					Request: &http.Request{
						Body: &errReader{},
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			parseErr: true,
		},
		{
			name: "parse string",
			args: args{
				field: req5,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: `{"x":1,"y":"2"}`,
		},
		{
			name: "parse string failed",
			args: args{
				field: req5,
				context: &gin.Context{
					Request: &http.Request{
						Body: &errReader{},
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			parseErr: true,
		},
		{
			name: "parse any",
			args: args{
				field: req6,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: func(m any) bool {
				m2, ok := m.(map[string]any)
				if !ok {
					return false
				}
				return m2["x"] == float64(1) && m2["y"] == "2"
			},
		},
		{
			name: "parse any failed",
			args: args{
				field: req6,
				context: &gin.Context{
					Request: &http.Request{
						Body: &errReader{},
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			parseErr: true,
		},
		{
			name: "parse []any",
			args: args{
				field: req7,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`[1,2,3]`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: func(m any) bool {
				m2, ok := m.([]any)
				if !ok {
					return false
				}
				return len(m2) == 3 && m2[0] == float64(1) && m2[1] == float64(2) && m2[2] == float64(3)
			},
		},
		{
			name: "parse map[string]any",
			args: args{
				field: req8,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			want: func(m any) bool {
				m2, ok := m.(map[string]any)
				if !ok {
					return false
				}
				return m2["x"] == float64(1) && m2["y"] == "2"
			},
		}, {
			name: "parse *string",
			args: args{
				field: req9,
				context: &gin.Context{
					Request: &http.Request{
						Body: io.NopCloser(bytes.NewBufferString(`{"x":1,"y":"2"}`)),
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			buildErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bodyNameParser{}
			parse, err := b.BuildParser(nil, tt.args.field)
			if (err != nil) != tt.buildErr {
				t.Errorf("BuildParser() error = %v, wantErr %v", err, tt.buildErr)
				return
			}
			if err == nil {
				value, err := parse(tt.args.context)
				if (err != nil) != tt.parseErr {
					t.Errorf("BuildParser() error = %v, wantErr %v", err, tt.parseErr)
					return
				}
				if err == nil {

					if f, ok := tt.want.(func(any) bool); ok {
						if !f(value.Interface()) {
							t.Errorf("BuildParser() value = %v, want %v", value.Interface(), tt.want)
						}
						return
					}

					if tt.want != nil && !reflect.DeepEqual(value.Interface(), tt.want) {
						t.Errorf("BuildParser() value = %v, want %v", value.Interface(), tt.want)
					}
				}
			}
		})
	}
}

func Test_bodyNameParser_Name(t *testing.T) {

	tests := []struct {
		name string
		want string
	}{
		{
			name: "test",
			want: "body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bodyNameParser{}
			assert.Equalf(t, tt.want, b.Name(), "Name()")
		})
	}
}
