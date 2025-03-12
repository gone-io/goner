# HTTP 注入说明

## HTTP 依赖注入标签的格式

```
${attributeName} ${attributeType} gone:"http,${kind}=${key}"
```

举例：
```go
router.GET("/search", function(in struct{
    selects []int `gone:"http,query=select"`
}){
    //注入值in.selects为`[]int{1,2,3}`
    fmt.Printf("%v", in.selects)
})
```
上面例子中，
- `selects`为属性名（attributeName）；
- `[]int`为属性类型（attributeType）；
- `query`为注入类型（kind）；
- `select`为注入键值（key）。

## 支持注入的类型和响应标签

| 名称                | 属性类型`${attributeType}`                                       | 注入类型`${kind}` | 注入键值`${key}` | 说明                                                                                                                                                                                   |
| ------------------- | ---------------------------------------------------------------- | :---------------: | :--------------: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **上下文注入**      | `gone.Context`                                                   |         /         |        /         | （不推荐）注入gin请求上下文对象，不需要类型`${kind}`和键值`${key}`。                                                                                                                   |
| **上下文注入**      | `*gone.Context`                                                  |         /         |        /         | （推荐）注入gin请求上下文指针，不需要类型`${kind}`和键值`${key}`。                                                                                                                     |
| **请求注入**        | `http.Request`                                                   |         /         |        /         | 不推荐）注入http.Request对象，不需要类型`${kind}`和键值`${key}`。                                                                                                                      |
| **请求注入**        | `*http.Request`                                                  |         /         |        /         | （推荐）注入http.Request指针，不需要类型`${kind}`和键值`${key}`。                                                                                                                      |
| **地址注入**        | `url.URL`                                                        |         /         |        /         | （不推荐）注入url.URL，不需要类型`${kind}`和键值`${key}`。                                                                                                                             |
| **地址注入**        | `*url.URL`                                                       |         /         |        /         | （推荐）注入url.URL指针，不需要类型`${kind}`和键值`${key}`。                                                                                                                           |
| **请求头注入**      | `http.Header`                                                    |         /         |        /         | （推荐）注入http.Header（请求头），不需要类型`${kind}`和键值`${key}`。                                                                                                                 |
| **响应注入**        | `gone.ResponseWriter`                                            |         /         |        /         | 注入gin.ResponseWriter（用于直接写入响应数据），不需要类型`${kind}`和键值`${key}`。                                                                                                    |
| **Body注入**        | 结构体、结构体指针                                               |      `body`       |        /         | **body注入**；将请求body解析后注入到属性，注入类型为 `body`，不需要“注入键值`${key}`”；框架根据`Content-Type`自动判定是json还是xml等格式；每个请求处理函数只允许存在一个**body注入**。 |
| **请求头单值注入**  | number \| string                                                 |      header       |   缺省取字段名   | 以键值`${key}`为`key`获取请求头，属性类型支持 简单类型<sub>[1]</sub>，解析不了会返回参数错误                                                                        |
| **URL路径参数注入** | number \| string                                                 |       param       |   缺省取字段名   | 以“注入键值`${key}`”为`key`调用函数`ctx.Param(key)`获取Url中定义的参数值，属性类型支持 简单类型<sub>[1]</sub>，解析不了会返回参数错误                               |
| **Query参数注入**   | number \| string \| []number \| []string \| 结构体 \| 结构体指针 |       query       |   缺省取字段名   | 以“注入键值`${key}`”为`key`调用函数`ctx.Query(key)`获取Query中的参数，属性类型支持 简单类型<sub>[1]</sub>，**支持简单类型的数组**，支持结构体和结构体指针，解析不了会返回参数错误       |
| **Cookie注入**      | number \| string                                                 |      cookie       |   缺省取字段名   | 以“注入键值`${key}`”为`key`调用函数`ctx.Context.Cookie(key)`获取Cookie的值，属性类型支持 简单类型<sub>[1]</sub>，解析不了会返回参数错误                             |


## Query参数注入

### 属性类型为简单类型<sub>[1]</sub>
解析不了会返回参数错误。

```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				yourName string `gone:"http,query=name"` //注册请求query中的name参数
				name string `gone:"http,query"` //注册请求query中的name参数；不指定参数名，则取属性名作为参数名
                age int `gone:"http,query=age"` //int类型
			}) string {
				return fmt.Sprintf("hello, %s, your name is %s", in.yourName, in.name)
			},
		)
```
### 属性类型为简单类型的数组
解析不了会返回参数错误。
下面代码，query为`?keyword=gone&keyword=is&keyword=best`，`in.keywords`的值将会为 `[]string{"gone","is","best"}`。

```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				keywords []string `gone:"http,query=keyword"` //简单类型的query数组参数注入
			}) string {
				return fmt.Sprintf("hello, keywords is [%v]", in.keywords)
			},
		)
```

### 属性类型为结构体或者结构体指针
这种类型key无需指定；假设query为`?page=1&pageSize=20&keyword=gone&keyword=is&keyword=best`，`in.req`的值将会为 `{1,20,[]string{"gone","is","best"}}`；注意结构体中可以使用`form`标签进行属性映射。

解析不了会返回参数错误。
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

## URL路径参数注入
URL 路径参数，是指定义在URL路由中的参数，注入属性的类型支持`string`和`int`,`uint`,`float64`等数字类型，解析不了会返回参数错误。如下：
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello/:name", //url中定义参数名为name
			func(in struct {
				name string `gone:"http,param"` //不指定参数名，使用属性名作为参数名
				name2 string `gone:"http,param=name"` //使用key指定参数名
			}) string {
				return "hello, " + in.name
			},
		)
```

## Body注入
Body注入，是指读取HTTP请求正文内容，解析成结构体，注入属性的类型支持结构体、结构体指针，解析不了会返回参数错误。

支持多种ContentType：json、xml、form-data、form-urlencoded等，不传ContentType时，默认为application/x-www-form-urlencoded。

结构体可以使用`form`标签进行form-data、form-urlencoded的属性映射；`xml`标签进行xml的属性映射；`json`标签进行json的属性映射。

具体规则可以参考：[gin#Model binding and validation](https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation)。

举例如下：
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
				req Req `gone:"http,body"`  //注意：body只能被注入一次，因为 writer被读取后就变成空了
				// req2 *Req `gone:"http,body"`
			}) string {
				fmt.Println(in.req)
				return fmt.Sprintf("hello, keywords is [%v]", in.req.Keywords)
			},
		)
```

## 请求头注入
HTTP请求头注入，用于获取某个请求头信息，注入属性的类型支持`string`和`int`,`uint`,`float64`等数字类型，解析不了会返回参数错误。
比如下面代码，可以用于读取请求头中的`Content-Type`信息。
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				contentType string `gone:"http,header"` //不指定参数名，使用属性名作为参数名
				contentType2 string `gone:"http,header=contentType"` //使用key指定参数名
			}) string {
				return "hello, contentType = " + in.contentType
			},
		)
```

## Cookie注入
Cookie注入，用于获取某个cookie信息，注入属性的类型支持`string`和`int`,`uint`,`float64`等数字类型，解析不了会返回参数错误。
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				token string `gone:"http,cookie"` //不指定参数名，使用属性名作为参数名
				token2 string `gone:"http,header=token"` //使用key指定参数名
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```

## 高级
另外，我们还支持几种特殊结构体（或结构体指针、接口、map）的注入，由于golang的“值拷贝”推荐使用指针注入，这些结构体代表了HTTP请求、响应、上下文等，这些结构体的注入不需要指定`kind`和`key`。

### URL结构体注入
支持属性类型为 `url.URL` 或者 `*url.URL`，该类型定义在`net/url`包中，代表了HTTP请求的URL。
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				url *url.URL `gone:"http"` //使用结构体指针
				url2 url.URL `gone:"http"` //使用结构体
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```

### 请求头注入
支持属性类型为 `http.Header`，该类型定义在`net/http`包中，代表了HTTP请求的Header。
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

### 上下文结构体注入
支持属性类型为 `gin.Content` 或者 `*gin.Content`，该类型定义在`github.com/gin-gonic/gin`包中，代表了HTTP请求的上下文。
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				context *gin.Content `gone:"http"` //使用结构体指针
				context2 gin.Content `gone:"http"` //使用结构体
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```

### 请求结构体注入
支持属性类型为 `http.Request` 或者 `*http.Request`，该类型定义在`net/http`包中，代表了HTTP请求信息。
```go
	ctr.rootRouter.
		Group("/demo").
		POST(
			"/hello",
			func(in struct {
				request *http.Request `gone:"http"` //使用结构体指针
				request2 http.Request `gone:"http"` //使用结构体
			}) string {
				return "hello, your token in cookie is" + in.token
			},
		)
```
### 请求响应接口注入

支持属性类型为 `gin.ResponseWriter`，该类型定义在`github.com/gin-gonic/gin`包中，代表了HTTP响应信息，可以使用该接口响应请求信息。
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


## 备注
[1]. 简单类型指 字符串、布尔类型 和 数字类型，其中数字类型包括：
- 整数类型：int、uint、int8、uint8、int16、uint16、int32、uint32、int64、uint64
- 非负整数类型：uint、uint8、uint16、uint32、uint64
- 浮点类型：float32、float64