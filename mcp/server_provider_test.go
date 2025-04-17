package goneMcp

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

type iTool struct {
	gone.Flag
}

func (i iTool) Define() Tool {
	return mcp.NewTool("test", mcp.WithDescription("test tool"))
}

func (i iTool) Process() func(ctx context.Context, request CallToolRequest) (*CallToolResult, error) {
	return func(ctx context.Context, request CallToolRequest) (*CallToolResult, error) {
		return mcp.NewToolResultText("test"), nil
	}
}

var _ ITool = (*iTool)(nil)

type iPrompt struct {
	gone.Flag
}

func (i iPrompt) Define() Prompt {
	return mcp.NewPrompt("test", mcp.WithPromptDescription("test prompt"))
}

func (i iPrompt) Process() func(ctx context.Context, request GetPromptRequest) (*GetPromptResult, error) {
	return func(ctx context.Context, request GetPromptRequest) (*GetPromptResult, error) {
		return mcp.NewGetPromptResult("test prompt", []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent("test prompt"),
			),
		}), nil
	}
}

var _ IPrompt = (*iPrompt)(nil)

type iResource struct {
	gone.Flag
}

func (i iResource) Define() Resource {
	return mcp.NewResource("test://test", "test")
}

func (i iResource) Process() func(ctx context.Context, request ReadResourceRequest) ([]ResourceContents, error) {
	return func(ctx context.Context, request ReadResourceRequest) ([]ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "test://test",
				MIMEType: "text/markdown",
				Text:     "test content",
			},
		}, nil
	}
}

var _ IResource = (*iResource)(nil)

func Test_serverProvider_Init(t *testing.T) {
	t.Run("goner define", func(t *testing.T) {
		gone.
			NewApp(serverLoad).
			Load(&iTool{}).
			Load(&iPrompt{}).
			Load(&iResource{}).
			Run(func(in struct {
				s1 *server.MCPServer `gone:"*"`
				s2 *server.MCPServer `gone:"*"`
			}) {
				assert.NotNil(t, in.s1)
				assert.Equal(t, in.s1, in.s2)
			})
	})

	t.Run("goner define and sse", func(t *testing.T) {
		_ = os.Setenv("GONE_MCP", `{"transportType":"sse"}`)
		defer func() {
			_ = os.Unsetenv("GONE_MCP")
		}()

		gone.
			NewApp(serverLoad).
			Load(&iTool{}).
			Load(&iPrompt{}).
			Load(&iResource{}).
			Run(func(in struct {
				s1 *server.MCPServer `gone:"*"`
				s2 *server.MCPServer `gone:"*"`
			}) {
				assert.NotNil(t, in.s1)
				assert.Equal(t, in.s1, in.s2)
			})
	})

	t.Run("goner define and sse and init error", func(t *testing.T) {
		_ = os.Setenv("GONE_MCP", `{"transportType":"sse"}`)
		_ = os.Setenv("GONE_MCP_SSE", `err`)
		defer func() {
			_ = os.Unsetenv("GONE_MCP")
		}()

		defer func() {
			a := recover()
			assert.NotNil(t, a)
		}()

		gone.
			NewApp(serverLoad).
			Load(&iTool{}).
			Load(&iPrompt{}).
			Load(&iResource{}).
			Run(func(in struct {
				s1 *server.MCPServer `gone:"*"`
				s2 *server.MCPServer `gone:"*"`
			}) {
				assert.NotNil(t, in.s1)
				assert.Equal(t, in.s1, in.s2)
			})
	})
}

func Test_serverProvider_Provide(t *testing.T) {
	tests := []struct {
		name       string
		loader     gone.LoadFunc
		before     func() (after func())
		tagConf    string
		errFn      func(err error)
		checkValue func(v *server.MCPServer)
	}{
		{
			name:    "customer prefix mcp server",
			tagConf: "custom",
		},
		{
			name: "customer prefix mcp server, get mcp config error",
			before: func() (after func()) {
				_ = os.Setenv("GONE_CUSTOM", `err`)

				return func() {
					_ = os.Unsetenv("GONE_CUSTOM")
				}
			},
			tagConf: "custom",
			errFn: func(err error) {
				assert.Contains(t, err.Error(), "invalid character")
				assert.Contains(t, err.Error(), "get mcp server config failed by key=custom")
			},
		},
		{
			name:   "default prefix and inject hooks",
			loader: g.NamedThirdComponentLoadFunc("mcp.hooks", &server.Hooks{}),
			checkValue: func(v *server.MCPServer) {
				elem := reflect.ValueOf(v).Elem()
				hook := elem.FieldByName("hooks")
				magic := gone.BlackMagic(hook)
				assert.NotNil(t, magic.Interface())
			},
		},
		{
			name: "default prefix and not inject hooks",
			checkValue: func(v *server.MCPServer) {
				elem := reflect.ValueOf(v).Elem()
				hook := elem.FieldByName("hooks")
				magic := gone.BlackMagic(hook)
				assert.Nil(t, magic.Interface())
			},
		},
		{
			name: "inject StdioContextFunc",
			loader: g.NamedThirdComponentLoadFunc("mcp.stdio.context", server.StdioContextFunc(func(ctx context.Context) context.Context {
				return context.WithValue(ctx, "test", "test")
			})),
		},
		{
			name: "unsupported TransportType",
			before: func() (after func()) {
				_ = os.Setenv("GONE_CUSTOM", `{"transportType":"test"}`)
				return func() {
					_ = os.Unsetenv("GONE_CUSTOM")
				}
			},
			tagConf: "custom",
			errFn: func(err error) {
				assert.Contains(t, err.Error(), "unsupported TransportType")
			},
		},
		{
			name: "sse get conf error",
			before: func() (after func()) {
				_ = os.Setenv("GONE_CUSTOM", `{"transportType":"sse"}`)
				_ = os.Setenv("GONE_CUSTOM_SSE", `err`)
				return func() {
					_ = os.Unsetenv("GONE_CUSTOM")
					_ = os.Unsetenv("GONE_CUSTOM_SSE")
				}
			},
			tagConf: "custom",
			errFn: func(err error) {
				assert.Contains(t, err.Error(), "invalid character")
				assert.Contains(t, err.Error(), "get mcp sse server config failed by key=custom.sse")
			},
		},
		{
			name: "sse inject contextFn",
			before: func() (after func()) {
				_ = os.Setenv("GONE_CUSTOM", `{"transportType":"sse"}`)
				return func() {
					_ = os.Unsetenv("GONE_CUSTOM")
				}
			},
			loader: g.NamedThirdComponentLoadFunc("custom.sse.context",
				server.SSEContextFunc(func(ctx context.Context, r *http.Request) context.Context {
					return context.WithValue(ctx, "test", "test")
				})),
			tagConf: "custom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				defer tt.before()()
			}
			gone.
				NewApp(serverLoad).
				Loads(func(loader gone.Loader) error {
					if tt.loader != nil {
						return tt.loader(loader)
					}
					return nil
				}).
				Run(func(p *serverProvider) {
					v, err := p.Provide(tt.tagConf)
					if tt.errFn != nil {
						tt.errFn(err)
					} else {
						assert.Nil(t, err)
						assert.NotNil(t, v)
						if tt.checkValue != nil {
							tt.checkValue(v)
						}
					}
				})
		})
	}
}

func TestConf_ToOptions(t *testing.T) {
	type fields struct {
		Name                   string
		Version                string
		WithRecovery           bool
		WithPromptCapabilities bool
		WithToolCapabilities   bool
		WithLogging            bool
		WithInstructions       string
		TransportType          string
	}
	tests := []struct {
		name   string
		fields fields
		want   []server.ServerOption
	}{
		{
			name: "mcp server options",
			fields: fields{
				Name:                   "test",
				Version:                "1.0.0",
				WithRecovery:           true,
				WithPromptCapabilities: true,
				WithToolCapabilities:   true,
				WithLogging:            true,
				WithInstructions:       "test",
				TransportType:          "sse",
			},
			want: []server.ServerOption{
				server.WithRecovery(),
				server.WithPromptCapabilities(true),
				server.WithToolCapabilities(true),
				server.WithLogging(),
				server.WithInstructions("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := Conf{
				Name:                   tt.fields.Name,
				Version:                tt.fields.Version,
				WithRecovery:           tt.fields.WithRecovery,
				WithPromptCapabilities: tt.fields.WithPromptCapabilities,
				WithToolCapabilities:   tt.fields.WithToolCapabilities,
				WithLogging:            tt.fields.WithLogging,
				WithInstructions:       tt.fields.WithInstructions,
				TransportType:          tt.fields.TransportType,
			}
			options := conf.ToOptions()

			assert.Equal(t, len(tt.want), len(options))
		})
	}
}

func TestSSEConfig_ToOptions(t *testing.T) {
	type fields struct {
		WithBaseURL                      string
		WithBasePath                     string
		WithMessageEndpoint              string
		WithUseFullURLForMessageEndpoint bool
		WithSSEEndpoint                  string
		WithKeepAliveInterval            time.Duration
		WithKeepAlive                    bool
		Address                          string
	}
	tests := []struct {
		name   string
		fields fields
		want   []server.SSEOption
	}{
		{
			name: "sse server options",
			fields: fields{
				WithBaseURL:                      "test",
				WithBasePath:                     "test",
				WithMessageEndpoint:              "test",
				WithUseFullURLForMessageEndpoint: true,
				WithSSEEndpoint:                  "test",
				WithKeepAliveInterval:            time.Second,
				WithKeepAlive:                    true,
			},
			want: []server.SSEOption{
				server.WithBaseURL("test"),
				server.WithBasePath("test"),
				server.WithMessageEndpoint("test"),
				server.WithUseFullURLForMessageEndpoint(true),
				server.WithSSEEndpoint("test"),
				server.WithKeepAliveInterval(time.Second),
				server.WithKeepAlive(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := SSEConfig{
				WithBaseURL:                      tt.fields.WithBaseURL,
				WithBasePath:                     tt.fields.WithBasePath,
				WithMessageEndpoint:              tt.fields.WithMessageEndpoint,
				WithUseFullURLForMessageEndpoint: tt.fields.WithUseFullURLForMessageEndpoint,
				WithSSEEndpoint:                  tt.fields.WithSSEEndpoint,
				WithKeepAliveInterval:            tt.fields.WithKeepAliveInterval,
				WithKeepAlive:                    tt.fields.WithKeepAlive,
				Address:                          tt.fields.Address,
			}
			assert.Equalf(t, len(tt.want), len(conf.ToOptions()), "ToOptions()")
		})
	}
}
