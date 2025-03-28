# HTTP Injection Documentation

## Format of HTTP Dependency Injection Tags

```
${attributeName} ${attributeType} gone:"http,${kind}=${key}"
```

Example:
```go
router.GET("/search", function(in struct{
    selects []int `gone:"http,query=select"`
}){
    //Injected value in.selects will be `[]int{1,2,3}`
    fmt.Printf("%v", in.selects)
})
```
In the above example:
- `selects` is the attribute name (attributeName);
- `[]int` is the attribute type (attributeType);
- `query` is the injection type (kind);
- `select` is the injection key (key).

## Supported Injection Types and Response Tags

| Name                | Attribute Type `${attributeType}`                                       | Injection Type `${kind}` | Injection Key `${key}` | Description                                                                                                                                                                                   |
| ------------------- | ---------------------------------------------------------------- | :---------------: | :--------------: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Context Injection**      | `gone.Context`                                                   |         /         |        /         | (Not recommended) Inject gin request context object, doesn't require type `${kind}` and key `${key}`.                                                                                                                   |
| **Context Injection**      | `*gone.Context`                                                  |         /         |        /         | (Recommended) Inject gin request context pointer, doesn't require type `${kind}` and key `${key}`.                                                                                                                     |
| **Request Injection**        | `http.Request`                                                   |         /         |        /         | (Not recommended) Inject http.Request object, doesn't require type `${kind}` and key `${key}`.                                                                                                                      |
| **Request Injection**        | `*http.Request`                                                  |         /         |        /         | (Recommended) Inject http.Request pointer, doesn't require type `${kind}` and key `${key}`.                                                                                                                      |
| **URL Injection**        | `url.URL`                                                        |         /         |        /         | (Not recommended) Inject url.URL, doesn't require type `${kind}` and key `${key}`.                                                                                                                             |
| **URL Injection**        | `*url.URL`                                                       |         /         |        /         | (Recommended) Inject url.URL pointer, doesn't require type `${kind}` and key `${key}`.                                                                                                                           |
| **Header Injection**      | `http.Header`                                                    |         /         |        /         | (Recommended) Inject http.Header (request headers), doesn't require type `${kind}` and key `${key}`.                                                                                                                 |
| **Response Injection**        | `gone.ResponseWriter`                                            |         /         |        /         | Inject gin.ResponseWriter (for directly writing response data), doesn't require type `${kind}` and key `${key}`.                                                                                                    |
| **Body Injection**        | Struct, struct pointer                                               |      `body`       |        /         | **Body injection**; Parse request body and inject into attribute, injection type is `body`, doesn't require "injection key `${key}`"; framework automatically determines format (json/xml/etc) based on `Content-Type`; Only one **body injection** allowed per request handler function. |
| **Single Header Value Injection**  | number \| string                                                 |      header       |   Default to field name   | Get header value using key `${key}` as `key`, attribute type supports simple types<sub>[1]</sub>, parsing failure will return parameter error                                                                        |
| **URL Path Parameter Injection** | number \| string                                                 |       param       |   Default to field name   | Use "injection key `${key}`" as `key` to call `ctx.Param(key)` and get parameter value from URL, attribute type supports simple types<sub>[1]</sub>, parsing failure will return parameter error                               |
| **Query Parameter Injection**   | number \| string \| []number \| []string \| struct \| struct pointer |       query       |   Default to field name   | Use "injection key `${key}`" as `key` to call `ctx.Query(key)` and get parameter from Query, attribute type supports simple types<sub>[1]</sub>, **supports arrays of simple types**, supports struct and struct pointer, parsing failure will return parameter error       |
| **Cookie Injection**      | number \| string                                                 |      cookie       |   Default to field name   | Use "injection key `${key}`" as `key` to call `ctx.Context.Cookie(key)` and get Cookie value, attribute type supports simple types<sub>[1]</sub>, parsing failure will return parameter error                             |


## Query Parameter Injection

### Attribute type is simple type<sub>[1]</sub>
Parsing failure will return parameter error.

```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				yourName string `gone:"http,query=name"` //Register name parameter from request query
				name string `gone:"http,query"` //Register name parameter from request query; when parameter name not specified, use attribute name as parameter name
                age int `gone:"http,query=age"` //int type
			}) string {
				return fmt.Sprintf("hello, %s, your name is %s", in.yourName, in.name)
			},
		)
```
### Attribute type is array of simple types
Parsing failure will return parameter error.
In below code, when query is `?keyword=gone&keyword=is&keyword=best`, `in.keywords` value will be `[]string{"gone","is","best"}`.

```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				keywords []string `gone:"http,query=keyword"` //Array of simple type query parameters injection
			}) string {
				return fmt.Sprintf("hello, keywords is [%v]", in.keywords)
			},
		)
```

### Attribute type is struct or struct pointer
This type doesn't require key specification; Assuming query is `?page=1&pageSize=20&keyword=gone&keyword=is&keyword=best`, `in.req` value will be `{1,20,[]string{"gone","is","best"}}`; Note struct can use `form` tag for property mapping.

Parsing failure will return parameter error.
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
			func(in struct {
				req Req `gone:"http,query"`
                req2 *Req `gone:"http,query"`
			}) string {
				fmt.Println(in.req)
				return fmt.Sprintf("hello, keywords is [%v]", in.req.Keywords)
			},
		)
```

## URL Path Parameter Injection
URL path parameters are parameters defined in URL routes, supported attribute types include `string` and numeric types like `int`,`uint`,`float64`, parsing failure will return parameter error. Example:
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello/:name", //Define parameter named name in URL
			func(in struct {
				name string `gone:"http,param"` //When parameter name not specified, use attribute name as parameter name
				name2 string `gone:"http,param=name"` //Use key to specify parameter name
			}) string {
				return "hello, " + in.name
			},
		)
```

## Body Injection
Body injection means reading HTTP request body content and parsing into struct, supported attribute types include struct and struct pointer, parsing failure will return parameter error.

Supports multiple ContentType: json, xml, form-data, form-urlencoded, etc. When ContentType not provided, defaults to application/x-www-form-urlencoded.

Struct can use `form` tag for form-data/form-urlencoded property mapping; `xml` tag for xml property mapping; `json` tag for json property mapping.

Specific rules can refer to: [gin#Model binding and validation](https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation).

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
			func(in struct {
				req Req `gone:"http,body"`  //Note: body can only be injected once, because writer becomes empty after being read
				// req2 *Req `gone:"http,body"`
			}) string {
				fmt.Println(in.req)
				return fmt.Sprintf("hello, keywords is [%v]", in.req.Keywords)
			},
		)
```

## Header Injection
HTTP header injection is used to get specific header information, supported attribute types include `string` and numeric types like `int`,`uint`,`float64`, parsing failure will return parameter error.
For example, below code can be used to read `Content-Type` information from request headers.
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				contentType string `gone:"http,header"` //When parameter name not specified, use attribute name as parameter name
				contentType2 string `gone:"http,header=contentType"` //Use key to specify parameter name
			}) string {
				return "hello, contentType = " + in.contentType
			},
		)
```

## Cookie Injection
Cookie injection is used to get specific cookie information, supported attribute types include `string` and numeric types like `int`,`uint`,`float64`, parsing failure will return parameter error.
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				token string `gone:"http,cookie"` //When parameter name not specified, use attribute name as parameter name
				token2 string `gone:"http,header=token"` //Use key to specify parameter name
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```

## Advanced
Additionally, we support injection of several special struct types (or struct pointers, interfaces, map). Due to golang's "value copy" mechanism, pointer injection is recommended. These structs represent HTTP request, response, context, etc., and their injection doesn't require specifying `kind` and `key`.

### URL Struct Injection
Supports attribute types `url.URL` or `*url.URL`, this type is defined in `net/url` package and represents HTTP request URL.
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				url *url.URL `gone:"http"` //Using struct pointer
				url2 url.URL `gone:"http"` //Using struct
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```

### Header Injection
Supports attribute type `http.Header`, this type is defined in `net/http` package and represents HTTP request Header.
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				header http.Header `gone:"http"`
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```

### Context Struct Injection
Supports attribute types `gin.Content` or `*gin.Content`, this type is defined in `github.com/gin-gonic/gin` package and represents HTTP request context.
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				context *gin.Content `gone:"http"` //Using struct pointer
				context2 gin.Content `gone:"http"` //Using struct
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```

### Request Struct Injection
Supports attribute types `http.Request` or `*http.Request`, this type is defined in `net/http` package and represents HTTP request information.
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				request *http.Request `gone:"http"` //Using struct pointer
				request2 http.Request `gone:"http"` //Using struct
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```
### Request Response Interface Injection

Supports attribute type `gin.ResponseWriter`, this type is defined in `github.com/gin-gonic/gin` package and represents HTTP response information, can use this interface to respond request information.
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				writer gin.ResponseWriter `gone:"http"`
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```


## Notes
[1]. Simple types refer to string, boolean and numeric types, where numeric types include:
- Integer types: int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64
- Unsigned integer types: uint, uint8, uint16, uint32, uint64
- Floating point types: float32, float64