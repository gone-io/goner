# Gone-Gin Performance Benchmark Report

## Introduction

Gone-Gin is an extension of the [Gin](https://github.com/gin-gonic/gin) framework, providing more convenient dependency injection and parameter binding features. It allows developers to focus on business logic implementation without worrying about request parameter parsing and binding. This document will introduce the implementation principles of Gone-Gin and its performance comparison with native Gin.

## Implementation Principles

### Core Components

1. **HTTP Injector (HttpInjector)**
   - Responsible for injecting HTTP request parameters into handler function parameters
   - Supports parameter injection from different sources (Body, Header, Query, Param, Cookie)
   - Implements automatic parameter binding and type conversion through reflection

2. **Proxy**
   - Converts user-defined handler functions into Gin-compatible handler functions
   - Supports multiple function signatures for more flexible programming patterns
   - Provides unified response processing for handler function execution results

3. **Response Processor (Responer)**
   - Converts handler function return values into HTTP responses
   - Supports multiple return value types including structs, Maps, slices, errors, etc.
   - Provides unified response format to simplify error handling
   - Supports custom response wrapper functions for flexible formatting
   - Built-in handling mechanisms for business errors and system errors
   - Supports advanced features like streaming responses (SSE) and file downloads

### Workflow

1. User defines handler functions with `gone:"http,xxx"` tags to specify parameter sources
2. Gone-Gin converts user functions into Gin-compatible handler functions through proxy mechanism
3. When a request arrives, the HTTP injector parses the request and injects parameters into user functions
4. After user function execution completes, results are uniformly processed by the response processor (Responer) and returned to the client

### Response Processing Mechanism

The response processor (Responer) is one of Gone-Gin's core components, responsible for converting handler function return values into HTTP responses. Compared to native Gin which requires manual response construction, Gone-Gin's response processing mechanism is more intelligent and flexible:

1. **Automatic Type Recognition**:
   - Complex types like structs, Maps, slices, arrays are automatically converted to JSON responses
   - Basic types (strings, numbers, etc.) are automatically converted to text responses
   - Error types are automatically processed as error responses with appropriate status codes

2. **Unified Response Format**:
   - Default unified response format `{"code": 0, "msg": "", "data": ...}`
   - Supports disabling unified format through configuration to return raw data directly
   - Supports custom response wrapper functions for flexible formatting

3. **Error Handling Mechanism**:
   - Distinguishes between business errors (BusinessError) and system errors (InnerError)
   - Business errors are returned to clients with error codes, messages and related data
   - System errors are properly handled to avoid sensitive information leakage while logging detailed information

4. **Advanced Feature Support**:
   - Supports streaming responses (SSE), automatically handling Channel type return values
   - Supports file downloads, automatically handling io.Reader type return values
   - Supports multiple return value processing with automatic recognition and error handling

## Performance Testing

### Test Environment

```
goos: darwin
goarch: arm64
cpu: Apple M1 Pro
```

### Test Cases

We designed four benchmark test cases to compare Gone-Gin with native Gin performance:

1. **BenchmarkProcessRequestWithInject**: Processes complete HTTP requests using Gone-Gin
2. **BenchmarkProcessRequestWithOriGin**: Processes complete HTTP requests using native Gin
3. **BenchmarkProxyGinHandlerFunc**: Tests performance of Gone-Gin proxy mechanism
4. **BenchmarkCallOriGinHandlerFunc**: Tests performance of native Gin handler functions

#### Test Code Examples

##### Request Struct Definition

```go
type Req struct {
	A int    `json:"a,omitempty"`
	B int    `json:"b,omitempty"`
	C int    `json:"c,omitempty"`
	D int    `json:"d,omitempty"`
	E string `json:"e,omitempty"`
	F string `json:"f,omitempty"`
}
```

##### Gone-Gin Handler Function

```go
// Gone-gin http handler function
func (c *ctr) httpHandler(in struct {
	req *Req `gone:"http,body"`
}) string {
	return "ok"
}
```

##### Native Gin Handler Function

```go
// Native Gin handler function
func originHandler(c *gin.Context) {
	var req Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}
```

##### Benchmark Functions

```go
// BenchmarkProcessRequestWithInject tests processing requests with gone-gin
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

// BenchmarkProcessRequestWithOriGin tests processing requests with native gin
func BenchmarkProcessRequestWithOriGin(b *testing.B) {
	engine := gin.New()
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
```

#### HTTP Injector Core Code

```go
// HTTP injector core implementation - handling Body parameter injection
func (s *httpInjector) injectBody(kind, key string, field reflect.StructField) (fn BindFieldFunc, err error) {
	if s.isInjectedBody {
		return nil, cannotInjectBodyMoreThanOnce(field.Name)
	}

	t := field.Type
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		return func(ctx *gin.Context, structVale reflect.Value) error {
			v := fieldByIndexFromStructValue(structVale, field.Index, field.IsExported(), field.Type)
			body := reflect.New(t).Interface()

			if err := ctx.ShouldBind(body); err != nil {
				return NewParameterError(err.Error())
			}
			v.Set(reflect.ValueOf(body).Elem())
			return nil
		}, nil
	case reflect.Pointer:
		return func(ctx *gin.Context, structVale reflect.Value) error {
			v := fieldByIndexFromStructValue(structVale, field.Index, field.IsExported(), field.Type)
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			if err := ctx.ShouldBind(v.Interface()); err != nil {
				return NewParameterError(err.Error())
			}
			return nil
		}, nil
	// ... other type handling
	default:
		return nil, unsupportedAttributeType(field.Name)
	}
}
```

#### Proxy Mechanism Core Code

```go
// Proxy mechanism core implementation - building proxy function
func (p *proxy) buildProxyFn(x HandlerFunc, funcName string, last bool) gin.HandlerFunc {
	m := make(map[int]*bindStructFuncAndType)
	args, err := p.funcInjector.InjectFuncParameters(
		x,
		func(pt reflect.Type, i int, injected bool) any {
			switch pt {
			case ctxPointType, ctxType, goneContextPointType, goneContextType:
				return placeholder{
					Type: pt,
				}
			}
			p.injector.StartBindFuncs()
			return nil
		},
		func(pt reflect.Type, i int, injected bool) any {
			m[i] = &bindStructFuncAndType{
				Fn:   p.injector.BindFuncs(),
				Type: pt,
			}
			return nil
		},
	)

	// ... error handling

	fv := reflect.ValueOf(x)
	return func(context *gin.Context) {
		// ... performance statistics

		parameters := make([]reflect.Value, 0, len(args))
		for i, arg := range args {
			// ... parameter processing
			if f, ok := m[i]; ok {
				parameter, err := f.Fn(context, arg)
				if err != nil {
					p.responser.Failed(context, err)
					return
				}
				parameters = append(parameters, parameter)
			} else {
				parameters = append(parameters, arg)
			}
		}

		// Call user function
		values := fv.Call(parameters)

		// Process return values
		var results []any
		for i := 0; i < len(values); i++ {
			// ... return value processing
		}
		p.responser.ProcessResults(context, context.Writer, last, funcName, results...)
	}
}
```

### Test Results

```
BenchmarkProcessRequestWithInject-8   261738     4124 ns/op
BenchmarkProcessRequestWithOriGin-8   370796     3166 ns/op
BenchmarkProxyGinHandlerFunc-8        278386     3942 ns/op
BenchmarkCallOriGinHandlerFunc-8      387363     3016 ns/op
```

### Result Analysis

1. **Complete Request Processing Performance**:
   - Gone-Gin: 4124 ns/op
   - Native Gin: 3166 ns/op
   - Performance difference: ~30%

2. **Handler Function Performance**:
   - Gone-Gin proxy: 3942 ns/op
   - Native Gin function: 3016 ns/op
   - Performance difference: ~31%

The test results show that Gone-Gin has some performance overhead compared to native Gin, but within the same order of magnitude. This is mainly due to:

1. **Reflection Mechanism**: Gone-Gin uses reflection for parameter injection and type conversion, which introduces performance overhead
2. **Proxy Layer**: Gone-Gin adds a proxy layer to convert user functions, increasing function call chain length
3. **Additional Features**: Gone-Gin provides more features like automatic parameter binding and type conversion, which introduce additional performance overhead

## Optimization Suggestions

Although Gone-Gin has some performance overhead compared to native Gin, this difference is acceptable in most business scenarios as it provides better development experience and higher efficiency. For further performance optimization, consider:

1. **Reduce Reflection Usage**: Minimize reflection usage or use caching to reduce reflection overhead
2. **Optimize Proxy Layer**: Simplify proxy layer implementation to reduce function call chain length
3. **On-Demand Loading**: Load only necessary features based on actual requirements to reduce unnecessary overhead
4. **Use More Efficient Data Structures**: Use more efficient data structures for storing and processing data when handling large volumes of requests

## Conclusion

Gone-Gin significantly improves development efficiency by providing more convenient dependency injection and parameter binding features. Although it has some performance overhead compared to native Gin, this difference is acceptable in most business scenarios. When choosing between Gone-Gin and native Gin, consider actual business requirements and performance needs.

For scenarios pursuing ultimate performance, consider using native Gin; for scenarios emphasizing development efficiency and code maintainability, Gone-Gin is a better choice.