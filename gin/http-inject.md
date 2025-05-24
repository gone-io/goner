<p>
    English&nbsp ｜&nbsp <a href="./http-inject_CN.md">中文</a>
</p>

# HTTP Injection Guide

## Overview

The HTTP Injector (HttpInjector) in goner/gin is a powerful dependency injection system that automatically injects various parameters from HTTP requests into handler function parameters. Through reflection mechanisms and tag parsing, it implements type-safe parameter binding and automatic type conversion.

## HTTP Dependency Injection Tag Format

```
${attributeName} ${attributeType} gone:"http,${kind}=${key}"
```

Example:

```go
router.GET("/search", function(in struct{
    selects []int `gone:"http,query=select"`
}){
    //injected value in.selects will be `[]int{1,2,3}`
    fmt.Printf("%v", in.selects)
})
```

In the above example:

- `selects` is the attribute name (attributeName);
- `[]int` is the attribute type (attributeType);
- `query` is the injection type (kind);
- `select` is the injection key value (key).

## Core Components

### 1. HTTP Injector (HttpInjector)

- Responsible for injecting HTTP request parameters into handler function parameters
- Supports parameter injection from different sources (Body, Header, Query, Param, Cookie)
- Implements automatic parameter binding and type conversion through reflection

### 2. Delay Bind Injector (DelayBindInjector)

- Implements delayed binding mechanism, parsing parameters only when the function is called
- Provides high-performance parameter injection capability
- Supports complex parameter types and nested structures

### 3. Bind Executor (BindExecutor)

- Manages type parsers (TypeParser) and name parsers (NameParser)
- Selects appropriate parsers based on tag configuration
- Executes specific parameter binding logic

## Supported Injection Types and Response Tags

### Types Supported by TypeParser

These types don't require specifying `kind` and `key`, injection is done directly through type matching:

| Name | Attribute Type `${attributeType}` | Description |
|------|----------------------------------|:------------|
| **Context Injection** | `*gin.Context` | Injects gin request context pointer. Implemented through `ginContextTypeParser`. |
| **Request Injection** | `*http.Request` | Injects http.Request pointer. Implemented through `httpRequestTypeParser`. |
| **URL Injection** | `*url.URL` | Injects url.URL pointer. Implemented through `urlTypeParser`. |
| **Header Injection** | `http.Header` | Injects http.Header (request headers). Implemented through `httpHeaderTypeParser`. |
| **Response Injection** | `gin.ResponseWriter` | Injects gin.ResponseWriter (for direct response writing). Implemented through `responseTypeParser`. |

### Types Supported by NameParser

These types require specifying `kind`, and optionally `key`:

| Name | Attribute Type `${attributeType}` | Injection Type `${kind}` | Injection Key `${key}` | Description |
|------|----------------------------------|:----------------------|:-------------------|:-------------|
| **Body Injection** | struct, struct pointer, []byte, io.Reader, io.ReadCloser, any, Map, Slice, String | `body` | / | **Body injection**; injects parsed request body into attribute, injection type is `body`, no need for injection key `${key}`; framework automatically determines format (json, xml, etc.) based on `Content-Type`; only one **body injection** allowed per request handler. Implemented through `bodyNameParser`. |
| **Single Header Injection** | number \| string | header | defaults to field name | Gets request header with key value `${key}` as `key`, attribute type supports simple types<sub>[1]</sub>, returns parameter error if parsing fails. Implemented through `headerNameParser`. |
| **URL Path Parameter Injection** | number \| string | param | defaults to field name | Gets URL parameter value by calling `ctx.Param(key)` with injection key value `${key}` as `key`, attribute type supports simple types<sub>[1]</sub>, returns parameter error if parsing fails. Implemented through `paramNameParser`. |
| **Query Parameter Injection** | number \| string \| []number \| []string \| struct \| struct pointer | query | defaults to field name | Gets query parameter with injection key value `${key}` as `key`, attribute type supports simple types<sub>[1]</sub>, **supports arrays of simple types**, supports structs and struct pointers, returns parameter error if parsing fails. Implemented through `queryNameParser`. |
| **Cookie Injection** | number \| string | cookie | defaults to field name | Gets cookie value by calling `ctx.Cookie(key)` with injection key value `${key}` as `key`, attribute type supports simple types<sub>[1]</sub>, returns parameter error if parsing fails. Implemented through `cookieNameParser`. |

## Implementation Principles

### Workflow

1. **Loading Phase**: Load HTTP injector and related parsers through `LoadGinHttpInjector` function
2. **Preparation Phase**: `DelayBindInjector` analyzes function parameter structure and selects appropriate parsers for each field
3. **Execution Phase**: When HTTP request arrives, extract and inject parameters using pre-compiled parsing functions

### Core Code Analysis

```go
// LoadGinHttpInjector loads the HTTP injector
func LoadGinHttpInjector(loader gone.Loader) error {
    loader.
        MustLoadX(injector.BuildLoad[*gin.Context](IdHttpInjector)).
        MustLoadX(parser.Load)
    return nil
}
```

This function completes loading of two key components:

1. `injector.BuildLoad[*gin.Context]` - Creates delayed binding injector for gin.Context
2. `parser.Load` - Loads all type parsers and name parsers

### Parser Architecture

#### Type Parsers (TypeParser)

- `ginContextTypeParser` - Handles `*gin.Context` type
- `httpRequestTypeParser` - Handles `*http.Request` type
- `httpHeaderTypeParser` - Handles `http.Header` type
- `urlTypeParser` - Handles `*url.URL` type
- `responseTypeParser` - Handles `gin.ResponseWriter` type

#### Name Parsers (NameParser)

- `bodyNameParser` - Handles `body` tag parameter injection
- `headerNameParser` - Handles `header` tag parameter injection
- `paramNameParser` - Handles `param` tag parameter injection
- `queryNameParser` - Handles `query` tag parameter injection
- `cookieNameParser` - Handles `cookie` tag parameter injection

## Query Parameter Injection

### Attribute Type as Simple Type<sub>[1]</sub>

Returns parameter error if parsing fails.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            yourName string `gone:"http,query=name"` //register name parameter from request query
            name string `gone:"http,query"`          //register name parameter from request query; if parameter name not specified, use attribute name
            age int `gone:"http,query=age"` //int type
        }) string {
            return fmt.Sprintf("hello, %s, your name is %s", in.yourName, in.name)
        },
    )
```

### Attribute Type as Array of Simple Types

Returns parameter error if parsing fails.
In the code below, for query `?keyword=gone&keyword=is&keyword=best`, `in.keywords` will be `[]string{"gone","is","best"}`.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            keywords []string `gone:"http,query=keyword"` //simple type query array parameter injection
        }) string {
            return fmt.Sprintf("hello, keywords is [%v]", in.keywords)
        },
    )
```

### Attribute Type as Struct or Struct Pointer

Key not required for this type; assuming query is `?page=1&pageSize=20&keyword=gone&keyword=is&keyword=best`, `in.req` value will be
`{1,20,[]string{"gone","is","best"}}`; note that `form` tags can be used for property mapping in the struct.

Returns parameter error if parsing fails.

```go
type Req struct {
    Page     string   `form:"page"`
    PageSize string   `form:"pageSize"`
    Keywords []string `form:"keywords"`
}

ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            req Req `gone:"http,query"`
            req2 *Req `gone:"http,query"`
        }) string {
            fmt.Println(in.req)
            return fmt.Sprintf("hello, keywords is [%v]", in.req.Keywords)
        },
    )
```

## URL Path Parameter Injection

URL path parameters are parameters defined in the URL route, injection attribute type supports `string` and numeric types like `int`, `uint`, `float64`, returns parameter error if parsing fails. Example:

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello/:name", //parameter name defined in url as name
        func (in struct {
            name string `gone:"http,param"`       //if parameter name not specified, use attribute name
            name2 string `gone:"http,param=name"` //use key to specify parameter name
        }) string {
            return "hello, " + in.name
        },
    )
```

## Body Injection

Body injection refers to reading HTTP request body content and parsing it into a struct, injection attribute type supports structs and struct pointers, returns parameter error if parsing fails.

Supports multiple ContentTypes: json, xml, form-data, form-urlencoded, etc. When ContentType not provided, defaults to application/x-www-form-urlencoded.

Structs can use `form` tags for form-data and form-urlencoded property mapping, `xml` tags for xml property mapping, and `json` tags for json property mapping.

For specific rules, refer to: [gin#Model binding and validation](https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation).

### Supported Body Types

Based on code analysis, `bodyNameParser` supports the following types:

1. **[]byte** - Direct raw byte data reading
2. **io.Reader / io.ReadCloser** - Returns request body's Reader interface
3. **struct/Map/Slice/any** - Uses gin's ShouldBind for automatic binding
4. **string** - Reads request body as string

Example:

```go
type Req struct {
    Page     string   `form:"page" json:"page,omitempty" xml:"page" binding:"required"`
    PageSize string   `form:"pageSize" json:"pageSize,omitempty" xml:"pageSize" binding:"required"`
    Keywords []string `form:"keywords" json:"keywords,omitempty" xml:"keywords" binding:"required"`
}

ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            req Req `gone:"http,body"` //note: body can only be injected once, as writer becomes empty after reading
            // req2 *Req `gone:"http,body"`
        }) string {
            fmt.Println(in.req)
            return fmt.Sprintf("hello, keywords is [%v]", in.req.Keywords)
        },
    )
```

## Header Injection

HTTP header injection is used to get specific header information, injection attribute type supports `string` and numeric types like `int`, `uint`, `float64`, returns parameter error if parsing fails.
For example, the code below can be used to read the `Content-Type` information from request headers.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            contentType string `gone:"http,header"`              //if parameter name not specified, use attribute name
            contentType2 string `gone:"http,header=contentType"` //use key to specify parameter name
        }) string {
            return "hello, contentType = " + in.contentType
        },
    )
```

## Cookie Injection

Cookie injection is used to get specific cookie information, injection attribute type supports `string` and numeric types like `int`, `uint`, `float64`, returns parameter error if parsing fails.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            token string `gone:"http,cookie"`        //if parameter name not specified, use attribute name
            token2 string `gone:"http,cookie=token"` //use key to specify parameter name
        }) string {
            return "hello, your token in cookie is" + in.token
        },
    )
```

## Advanced Usage

### Direct Type Parser Injection

For the following special types, direct injection is possible without specifying `kind` and `key`:

#### URL Struct Injection

Supports attribute type `*url.URL`, defined in the `net/url` package, represents the HTTP request URL.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            url *url.URL `gone:"http"` //use struct pointer
        }) string {
            return "hello, your url is " + url.String()
        },
    )
```

#### Complete Header Injection

Supports attribute type `http.Header`, defined in the `net/http` package, represents all HTTP request headers.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            header http.Header `gone:"http"`
        }) string {
            return "hello, your headers count is " + fmt.Sprintf("%d", len(header))
        },
    )
```

#### Context Struct Injection

Supports attribute type `*gin.Context`, defined in the `github.com/gin-gonic/gin` package, represents the HTTP request context.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            context *gin.Context `gone:"http"` //use struct pointer
        }) string {
            return "hello, your method is " + context.Request.Method
        },
    )
```

#### Request Struct Injection

Supports attribute type `*http.Request`, defined in the `net/http` package, represents HTTP request information.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            request *http.Request `gone:"http"` //use struct pointer
        }) string {
            return "hello, your method is " + request.Method
        },
    )
```

#### Response Interface Injection

Supports attribute type `gin.ResponseWriter`, defined in the `github.com/gin-gonic/gin` package, represents HTTP response information, can use this interface to respond to requests.

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            writer gin.ResponseWriter `gone:"http"`
        }) string {
        writer.Header().Set("Custom-Header", "custom-value")
            return "hello, custom header set"
        },
    )
```

### Types as Direct Function Parameters

The HTTP injector in goner/gin not only supports injecting parameters wrapped in structs but also supports using specific types directly as function parameters. This approach is more concise and suitable for scenarios requiring only a few parameters.

Types that can be used directly as function parameters are the same as those supported by TypeParser, including:
- `*gin.Context` - HTTP request context
- `*http.Request` - HTTP request object
- `*url.URL` - URL object
- `http.Header` - HTTP request headers
- `gin.ResponseWriter` - HTTP response writer

Example code:

```go
ctr.rootRouter.
    Group("/demo").
    POST("/users", func(ctx *gin.Context, req *http.Request, writer gin.ResponseWriter){
        // Use parameters directly, no need to extract from struct
        name := ctx.Query("name")
        writer.Header().Set("Content-Type", "application/json")
        writer.Write([]byte(fmt.Sprintf(`{"message":"Hello, %s"}`, name)))
    })
```

You can also mix struct parameters with direct parameters:

```go
ctr.rootRouter.
    Group("/demo").
    POST("/users", func(in struct {
        ID   int64  `gone:"http,param=id"`
        Name string `gone:"http,query=name"`
    }, writer gin.ResponseWriter){
        // Use both struct parameters and direct parameters
        writer.Header().Set("Content-Type", "application/json")
        writer.Write([]byte(fmt.Sprintf(`{"id":%d, "name":"%s"}`, in.ID, in.Name)))
    })
```

This approach makes the code more flexible, allowing you to choose the most suitable parameter passing method based on your needs.

### Custom Parameter Parsers

goner/gin allows developers to create custom parameter parsers to support more complex parameter injection scenarios.

#### Custom Type Parser

You can create a custom type parser by implementing the `injector.TypeParser[*gin.Context]` interface. This interface contains two methods:

- `Parse(context *gin.Context) (reflect.Value, error)`: Parses target type value from `gin.Context`.
- `Type() reflect.Type`: Returns the target type supported by the parser.

**Example: Custom Token Parser**

Suppose you need to parse a `Bearer Token` from the `Authorization` header and inject it into a custom `Token` type.

1.  **Define `Token` type:**

    ```go
    package main

    type Token string
    ```

2.  **Implement `TypeParser` interface:**

    ```go
    package main

    import (
        "reflect"
        "strings"

        "github.com/gin-gonic/gin"
        "github.com/gone-io/gone/v2"
        "github.com/gone-io/goner/gin/injector" // ensure correct import
    )

    // ensure tokenParser implements TypeParser interface
    var _ injector.TypeParser[*gin.Context] = (*tokenParser)(nil)

    type tokenParser struct {
        gone.Flag
    }

    func (t *tokenParser) Parse(context *gin.Context) (reflect.Value, error) {
        auth := context.GetHeader("Authorization")
        arr := strings.Split(auth, " ")
        if len(arr) == 2 && arr[0] == "Bearer" {
            token := Token(arr[1])
            return reflect.ValueOf(token), nil
        }
        return reflect.Value{}, gone.NewParameterError("invalid token") // use gone.NewParameterError to return error
    }

    func (t *tokenParser) Type() reflect.Type {
        return reflect.TypeOf(Token(""))
    }
    ```

3.  **Load custom parser:**

    ```go
    gone.Load(&tokenParser{})
    ```

4.  **Use in Handler:**

    Now you can directly inject the `Token` type in your Gin Handler.

    ```go
    package main

    import (
        "fmt"
        "github.com/gone-io/gone/v2"
        "github.com/gone-io/goner/gin"
    )

    type MyController struct {
        gone.Flag
        Router gin.IRouter `gone:"*"`
    }

    func (ctr *MyController) GetUser(in struct{
        token Token `gone:"http"`
    }) string { // directly inject Token type
        return fmt.Sprintf("User token: %s", in.token)
    }

    func (ctr *MyController) Mount() gin.MountError {
        ctr.Router.GET("/user", ctr.GetUser)
        return nil
    }
    ```

    When requesting the `/user` endpoint with the correct `Authorization: Bearer <your-token>` header, the `token` parameter will be automatically injected.

Through this approach, you can flexibly extend goner/gin's parameter injection capabilities to adapt to various complex business requirements.

#### Custom Name Parser

In addition to custom type parsers, goner/gin also allows you to create custom name parsers (NameParser). Name parsers handle parameters specified through the `gone:"http,kind=key"` tag. This provides finer-grained control, allowing you to implement custom parameter parsing logic based on specific `kind` and `key`.

You need to implement the `injector.NameParser[*gin.Context]` interface to create a custom name parser. Based on the example in `goner/gin/parser/name_query.go`, a typical name parser mainly includes the following methods:

-   `Name() string`: Returns the `kind` type handled by this name parser. For example, `query`, `header`, `param`, `cookie`, `body`, or your custom `kind`.
-   `BuildParser(keyMap map[string]string, field reflect.StructField) (func(context *gin.Context) (reflect.Value, error), error)`: This is the core build method. It's called during application initialization.
    -   `keyMap`: A key-value map parsed from the `gone:"http,kind=key,..."` tag. For example, for `gone:"http,query=userId,optional"`, `keyMap` might contain `{"query": "userId", "optional": ""}`. The value corresponding to `s.Name()` (i.e., `kind`) is the main `key`.
    -   `field`: The `reflect.StructField` information of the struct field being processed.
    -   This method needs to return a **parsing function** and an error. This parsing function `func(context *gin.Context) (reflect.Value, error)` will be actually executed when each HTTP request arrives, used to extract data from `*gin.Context`, convert and return `reflect.Value`.

This design pattern optimizes runtime performance by pre-building parsing logic during initialization.

**Example: Custom CSV Parser**

Suppose you need to get a comma-separated string from a query parameter and parse it into a string slice. For example, for request URL `/items?tags=go,gin,gone`, you want to inject it into a `Tags []string` field.

1.  **Implement `injector.NameParser[*gin.Context]` interface:**

    ```go
    package main

    import (
        "fmt"
        "reflect"
        "strings"

        "github.com/gin-gonic/gin"
        "github.com/gone-io/gone/v2"
        "github.com/gone-io/goner/gin/injector"
    )

    // ensure csvQueryParser implements NameParser interface
    var _ injector.NameParser[*gin.Context] = (*csvQueryParser)(nil)

    type csvQueryParser struct {
        gone.Flag
    }

    func (p *csvQueryParser) Name() string {
        return "csv" // custom kind as csv
    }

    func (p *csvQueryParser) BuildParser(keyMap map[string]string, field reflect.StructField) (func(context *gin.Context) (reflect.Value, error), error) {
        // get key corresponding to "csv" kind from keyMap
        // for example, for gone:"http,csv=my_tags", mainKey should be "my_tags"
        // for gone:"http,csv", mainKey might be empty, then can use field name field.Name
        mainKey := keyMap[p.Name()] 
        if mainKey == "" {
            mainKey = field.Name // if key not specified in tag, default to field name
        }

        // check if target field type is []string
        if field.Type.Kind() != reflect.Slice || field.Type.Elem().Kind() != reflect.String {
            return nil, fmt.Errorf("CSV parser: field '%s' must be of type []string, got %s", field.Name, field.Type.String())
        }

        // return actual parsing function
        return func(ctx *gin.Context) (reflect.Value, error) {
            paramValue := ctx.Query(mainKey)
            if paramValue == "" {
                // if field is required, can return gone.NewParameterError here
                // if allowed to be empty, return empty slice
                return reflect.ValueOf([]string{}), nil 
            }
            items := strings.Split(paramValue, ",")
            return reflect.ValueOf(items), nil
        }, nil
    }

    // factory function to be loaded by gone.Load
    func NewCsvQueryParser() gone.Goner {
        return &csvQueryParser{}
    }
    ```

2.  **Load custom parser:**

    You need to load your custom name parser into Gone's dependency injection container. This is typically done through `gone.Load()`, ensuring it's discovered by `GinHttpInjector`.
    `GinHttpInjector` will collect all components that implement the `injector.NameParser[*gin.Context]` interface.

    ```go
    // in your gone application startup logic
    gone.Load(NewCsvQueryParser())
    ```

3.  **Use in Handler:**

    Now you can use the `csv` kind in your Gin Handler's input struct.

    ```go
    package main

    import (
        "fmt"
        "github.com/gone-io/gone/v2"
        "github.com/gone-io/goner/gin"
    )

    type MyController struct {
        gone.Flag
        Router gin.IRouter `gone:"*"`
    }

    type ItemRequest struct {
        Tags []string `gone:"http,csv=item_tags"` // use custom csv parser
        Name string   `gone:"http,query=name"`    // use built-in query parser
    }

    func (ctr *MyController) CreateItem(in ItemRequest) string {
        return fmt.Sprintf("Item created with name '%s' and tags: %v", in.Name, in.Tags)
    }

    func (ctr *MyController) Mount() error {
        ctr.Router.POST("/items", ctr.CreateItem)
        return nil
    }
    ```

    When requesting `/items?name=MyItem&item_tags=urgent,important`:
    - `in.Name` will be injected as `"MyItem"` (through built-in `queryNameParser`)
    - `in.Tags` will be injected as `[]string{"urgent", "important"}` (through your `csvQueryParser`)

Through custom name parsers, you can greatly enhance goner/gin's flexibility and capability in handling HTTP request parameters, making it adapt to various complex API designs and data formats.

## Performance Optimization

### Delayed Binding Mechanism

The HTTP injector in goner/gin uses a delayed binding mechanism, pre-compiling parameter parsing functions during application startup and executing them directly during request handling, avoiding runtime reflection overhead and providing excellent performance.

### Type Safety

All parameter injections are type-safe, type mismatches can be detected at compile time, and detailed error information is provided at runtime to help debugging.

## Error Handling

When parameter parsing fails, the framework returns `gone.ParameterError` containing detailed error information to help developers quickly locate issues.

## Notes

[1]. Simple types refer to strings, boolean types, and numeric types, where numeric types include:

- Integer types: int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64
- Non-negative integer types: uint, uint8, uint16, uint32, uint64
- Floating-point types: float32, float64