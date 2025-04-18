package goneMcp

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type xTransport struct {
	gone.Flag
}

func (x xTransport) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) SendRequest(ctx context.Context, request transport.JSONRPCRequest) (*transport.JSONRPCResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) SendNotification(ctx context.Context, notification mcp.JSONRPCNotification) error {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) SetNotificationHandler(handler func(notification mcp.JSONRPCNotification)) {
	//TODO implement me
	panic("implement me")
}

func (x xTransport) Close() error {
	//TODO implement me
	panic("implement me")
}

var _ transport.Interface = (*xTransport)(nil)

func TestClientProviderProvide(t *testing.T) {
	t.Run("TestClientProviderProvide", func(t *testing.T) {
		gone.
			NewApp(ClientLoad).
			Load(&xTransport{}, gone.Name("x-transport")).
			Run(func(in struct {
				c1 *client.Client `gone:"*,type=stdio,param=go run ./testdata/stdio_server"`
				c2 *client.Client `gone:"*,type=sse,param=http://localhost:8082/sse"`
				c3 *client.Client `gone:"*,type=transport,param=x-transport"`
			}, in2 struct {
				c1 *client.Client `gone:"*,type=stdio,param=go run ./testdata/stdio_server"`
			}) {

			})
	})
	t.Run("stdio config with configKey", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			_ = os.Setenv("GONE_MCP_X_CLIENT", `{"command":"go", "args":["run", "./testdata/stdio_server"]}`)
			defer func() {
				_ = os.Unsetenv("GONE_MCP_X_CLIENT")
			}()

			gone.
				NewApp(ClientLoad).
				Load(&xTransport{}, gone.Name("x-transport")).
				Run(func(in struct {
					c1 *client.Client `gone:"*,type=stdio,configKey=mcp.x.client"`
				}) {

				})

		})
		t.Run("read config error", func(t *testing.T) {
			_ = os.Setenv("GONE_MCP_ERROR_CLIENT", `{"command":"go", "args":"err"}`)
			defer func() {
				_ = os.Unsetenv("GONE_MCP_X_CLIENT")
			}()

			gone.
				NewApp(ClientLoad).
				Load(&xTransport{}, gone.Name("x-transport")).
				Run(func(in struct {
					provider gone.Provider[*client.Client] `gone:"*"`
				}) {

					provide, err := in.provider.Provide("type=stdio,configKey=mcp.error.client")
					assert.NotNil(t, err)
					assert.Nil(t, provide)
					assert.Contains(t, err.Error(), "get mcp client config failed by key=mcp.error.client")
				})
		})
	})

	t.Run("sse read config with configKey", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			_ = os.Setenv("GONE_MCP_X_CLIENT", `{"baseUrl":"http://localhost:8082/sse", "header":{"x-test":"---test---"}}`)
			defer func() {
				_ = os.Unsetenv("GONE_MCP_X_CLIENT")
			}()

			gone.
				NewApp(ClientLoad).
				Load(&xTransport{}, gone.Name("x-transport")).
				Run(func(in struct {
					c1 *client.Client `gone:"*,type=sse,configKey=mcp.x.client"`
				}) {

				})

		})
		t.Run("read config error", func(t *testing.T) {
			_ = os.Setenv("GONE_MCP_ERROR_CLIENT", `err`)
			defer func() {
				_ = os.Unsetenv("GONE_MCP_X_CLIENT")
			}()

			gone.
				NewApp(ClientLoad).
				Load(&xTransport{}, gone.Name("x-transport")).
				Run(func(in struct {
					provider gone.Provider[*client.Client] `gone:"*"`
				}) {

					provide, err := in.provider.Provide("type=sse,configKey=mcp.error.client")
					assert.NotNil(t, err)
					assert.Nil(t, provide)
					assert.Contains(t, err.Error(), "get mcp client config failed by key=mcp.error.client")
				})
		})
	})

	t.Run("TestClientProviderProvideError", func(t *testing.T) {
		gone.
			NewApp(ClientLoad).
			Run(func(in struct {
				provider gone.Provider[*client.Client] `gone:"*"`
			}) {
				assert.NotNil(t, in.provider)
				provide, err := in.provider.Provide("")
				assert.NotNil(t, err)
				assert.Nil(t, provide)
				assert.Contains(t, err.Error(), "support type")

				c, err := in.provider.Provide("type=stdio")
				assert.NotNil(t, err)
				assert.Nil(t, c)
				assert.Contains(t, err.Error(), "create mcp client failed by type=stdio")

			})
	})

	t.Run("get injected transport failed", func(t *testing.T) {
		gone.
			NewApp(ClientLoad).
			Run(func(in struct {
				provider gone.Provider[*client.Client] `gone:"*"`
			}) {

				c, err := in.provider.Provide("type=transport,param=xxx")
				assert.NotNil(t, err)
				assert.Nil(t, c)
				assert.Contains(t, err.Error(), "can not found the transport by name=xxx, please load it first")

			})
	})

}
