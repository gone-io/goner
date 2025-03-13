package gin_test

import (
	origin "github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/tracer"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func init() {
	gone.GetDefaultLogger().SetLevel(gone.ErrorLevel)
}

type Req struct {
	A int    `json:"a,omitempty"`
	B int    `json:"b,omitempty"`
	C int    `json:"c,omitempty"`
	D int    `json:"d,omitempty"`
	E string `json:"e,omitempty"`
	F string `json:"f,omitempty"`
}

type ctr struct {
	gone.Flag
	gin.IRouter `gone:"*"`
}

func (c *ctr) Mount() gin.MountError {
	c.POST("/api/test", c.httpHandler)
	return nil
}

// Gone-gin 的 http 处理函数
func (c *ctr) httpHandler(in struct {
	req *Req `gone:"http,body"`
}) string {
	return "ok"
}

// 原生 gin 处理函数
func originHandler(c *origin.Context) {
	var req Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

// 构建请求和响应
func buildRequest() *http.Request {
	reqBody := `{"a": 1, "b": 2, "c": 3, "d": 4, "e": "test", "f": "example"}`
	request := httptest.NewRequest(http.MethodPost, "/api/test", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func buildResponse() http.ResponseWriter {
	return httptest.NewRecorder()
}

func buildGinContext(engine *origin.Engine) *origin.Context {
	ctx := origin.CreateTestContextOnly(buildResponse(), engine)
	ctx.Request = buildRequest()
	return ctx
}

// BenchmarkProcessRequestWithInject 测试 使用 gone-gin 处理请求
func BenchmarkProcessRequestWithInject(b *testing.B) {
	_ = os.Setenv("GONE_SERVER_SYS-MIDDLEWARE_DISABLE", "true")
	_ = os.Setenv("GONE_SERVER_RETURN_WRAPPED-DATA", "false")

	gone.
		NewApp(gin.Load, tracer.Load).
		Load(&ctr{}).
		Run(func(httpHandler http.Handler) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				request := buildRequest()
				response := buildResponse()
				b.StartTimer()
				httpHandler.ServeHTTP(response, request)
			}
		})
}

// BenchmarkProcessRequestWithOriGin 测试 使用 原生 gin 处理请求
func BenchmarkProcessRequestWithOriGin(b *testing.B) {
	engine := origin.New()
	engine.
		POST("/api/test", originHandler)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		request := buildRequest()
		response := buildResponse()
		b.StartTimer()
		engine.ServeHTTP(response, request)
	}
}

// BenchmarkProxyGinHandlerFunc 测试 调用 使用gone-gin生成的代理 HandlerFunc
func BenchmarkProxyGinHandlerFunc(b *testing.B) {
	_ = os.Setenv("GONE_SERVER_SYS-MIDDLEWARE_DISABLE", "true")
	_ = os.Setenv("GONE_SERVER_RETURN_WRAPPED-DATA", "false")
	engine := origin.New()
	gone.
		NewApp(gin.Load, tracer.Load).
		Load(&ctr{}).
		Run(func(proxy gin.HandleProxyToGin, ctr *ctr) {
			handle := proxy.Proxy(ctr.httpHandler)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				context := buildGinContext(engine)
				b.StartTimer()
				handle[0](context)
			}
		})
}

// BenchmarkCallOriGinHandlerFunc 测试 调用 编写的gin原生 HandlerFunc 效率
func BenchmarkCallOriGinHandlerFunc(b *testing.B) {
	engine := origin.New()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		context := buildGinContext(engine)
		b.StartTimer()
		originHandler(context)
	}
}
