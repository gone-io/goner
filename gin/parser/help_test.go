package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBuildParser(t *testing.T) {
	type X struct {
		Name string
	}

	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name     string
		args     args
		parseStr string

		wantErr   assert.ErrorAssertionFunc
		parseErr  assert.ErrorAssertionFunc
		parseWant any
	}{
		{
			name: "pt is string",
			args: args{
				t: reflect.TypeOf(""),
			},
			parseStr:  "hello",
			wantErr:   assert.NoError,
			parseErr:  assert.NoError,
			parseWant: "hello",
		},
		{
			name: "pt is struct",
			args: args{
				t: reflect.TypeOf(X{}),
			},
			parseStr:  "{\"Name\":\"hello\"}",
			wantErr:   assert.NoError,
			parseErr:  assert.NoError,
			parseWant: X{Name: "hello"},
		},
		{
			name: "pt is map && Unmarshal err",
			args: args{
				t: reflect.TypeOf(map[string]string{}),
			},
			parseStr: "{\"Name\":\"hello\"",
			wantErr:  assert.NoError,
			parseErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name: "pt is number",
			args: args{
				t: reflect.TypeOf(int(1)),
			},
			parseStr:  "1",
			wantErr:   assert.NoError,
			parseErr:  assert.NoError,
			parseWant: 1,
		},
		{
			name: "pt is number and error",
			args: args{
				t: reflect.TypeOf(int(1)),
			},
			parseStr: "x1",
			wantErr:  assert.NoError,
			parseErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Error(t, err, i...)
			},
		},
		{
			name: "pt is unsupported",
			args: args{
				t: reflect.TypeOf(make(chan int)),
			},
			parseStr: "1",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return !assert.Errorf(t, err, "unsupported type: %s", i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parse, err := BuildParser(tt.args.t)
			if !tt.wantErr(t, err, fmt.Sprintf("BuildParser(%v)", tt.args.t)) {
				return
			}
			value, err := parse(tt.parseStr)
			if !tt.parseErr(t, err, fmt.Sprintf("BuildParser(%v)", tt.args.t)) {
				return
			}
			if !reflect.DeepEqual(value.Interface(), tt.parseWant) {
				t.Errorf("BuildParser(%v) = %v, want %v", tt.args.t, value.Interface(), tt.parseWant)
			}
		})
	}
}
