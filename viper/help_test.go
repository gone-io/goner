package viper

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name   string
		before func() func()
		want   string
	}{
		{
			name: "default",
			before: func() func() {
				env := os.Getenv(EEnv)
				_ = os.Setenv(EEnv, "")
				return func() {
					_ = os.Setenv(EEnv, env)
				}
			},
			want: defaultEnv,
		},
		{
			name: "from env",
			before: func() func() {
				env := os.Getenv(EEnv)
				_ = os.Setenv(EEnv, "test")
				return func() {
					_ = os.Setenv(EEnv, env)
				}
			},
			want: "test",
		},
		{
			name: "from flag",
			before: func() func() {
				env := *envFlag
				*envFlag = "test"
				return func() {
					*envFlag = env
				}
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := tt.before()
			defer after()
			assert.Equalf(t, tt.want, GetEnv(), "GetEnv()")
		})
	}
}

func TestGetConfDir(t *testing.T) {
	tests := []struct {
		name   string
		before func() func()
		want   string
	}{
		{
			name: "default",
			before: func() func() {
				env := os.Getenv(EConf)
				_ = os.Setenv(EConf, "")
				return func() {
					_ = os.Setenv(EConf, env)
				}
			},
			want: "",
		},
		{
			name: "from env",
			before: func() func() {
				env := os.Getenv(EConf)
				_ = os.Setenv(EConf, "test")
				return func() {
					_ = os.Setenv(EConf, env)
				}
			},
			want: "test",
		},
		{
			name: "from flag",
			before: func() func() {
				env := *confFlag
				*confFlag = "test"
				return func() {
					*confFlag = env
				}
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := tt.before()
			defer after()
			assert.Equalf(t, tt.want, GetConfDir(), "GetConfDir()")
		})
	}
}

func TestMustGetExecutableConfDir(t *testing.T) {
	dir := MustGetExecutableConfDir()
	assert.NotEmpty(t, dir)
	assert.True(t, strings.HasSuffix(dir, ConPath))
}

func TestMustGetWorkDir(t *testing.T) {
	dir := MustGetWorkDir()
	assert.NotEmpty(t, dir)
}

func Test_lookForModDir(t *testing.T) {
	dir := lookForModDir(MustGetWorkDir())
	assert.NotEmpty(t, dir)
}

func Test_getConfigPaths(t *testing.T) {
	type args struct {
		isInTestKit bool
	}
	tests := []struct {
		name   string
		before func() func()
		args   args
		want   int
	}{
		{
			name: "isInTestKit=true",
			args: args{
				isInTestKit: true,
			},
			want: 7,
		},
		{
			name: "isInTestKit=false",
			args: args{
				isInTestKit: false,
			},
			want: 4,
		},
		{
			name: "isInTestKit=false,confDir=/tmp/test",
			before: func() func() {
				confDir := *confFlag
				*confFlag = "/tmp/test"
				return func() {
					*confFlag = confDir
				}
			},
			args: args{
				isInTestKit: false,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				after := tt.before()
				defer after()
			}
			assert.Equalf(t, tt.want, len(getConfigPaths(tt.args.isInTestKit)), "getConfigPaths(%v)", tt.args.isInTestKit)
		})
	}
}

func Test_findConfigFiles(t *testing.T) {
	type args struct {
		env         string
		isInTestKit bool
		paths       []string
		fsys        afero.Fs
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "env is empty & isInTestKit=true",
			args: args{
				env:         "",
				isInTestKit: true,
				paths:       []string{"testdata/config"},
				fsys:        afero.NewOsFs(),
			},

			want: []string{
				"testdata/config/default.json",
				"testdata/config/default.yaml",
				"testdata/config/default.properties",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "env not empty & isInTestKit=true",
			args: args{
				env:         "local",
				isInTestKit: true,
				paths:       []string{"testdata/config"},
				fsys:        afero.NewOsFs(),
			},

			want: []string{
				"testdata/config/default.json",
				"testdata/config/default.yaml",
				"testdata/config/default.properties",
				"testdata/config/local.properties",
				"testdata/config/local_test.yml",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "env not empty & isInTestKit=false",
			args: args{
				env:         "local",
				isInTestKit: false,
				paths:       []string{"testdata/config"},
				fsys:        afero.NewOsFs(),
			},

			want: []string{
				"testdata/config/default.json",
				"testdata/config/default.yaml",
				"testdata/config/default.properties",
				"testdata/config/local.properties",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findConfigFiles(tt.args.env, tt.args.isInTestKit, tt.args.paths, tt.args.fsys)
			if !tt.wantErr(t, err, fmt.Sprintf("findConfigFiles(%v, %v, %v, %v)", tt.args.env, tt.args.isInTestKit, tt.args.paths, tt.args.fsys)) {
				return
			}
			assert.Equalf(t, tt.want, got, "findConfigFiles(%v, %v, %v, %v)", tt.args.env, tt.args.isInTestKit, tt.args.paths, tt.args.fsys)
		})
	}
}

func Test_getConfigFiles(t *testing.T) {
	files, err := getConfigFiles(true, afero.NewOsFs())
	assert.NoError(t, err)
	assert.NotEmpty(t, files)
}

func Test_fileExt(t *testing.T) {
	ext := fileExt(".txt")
	assert.Equal(t, "txt", ext)
}
