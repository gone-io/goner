<p>
    <a href="http-inject.md">English</a>&nbsp ｜&nbsp 中文
</p>

# HTTP 注入说明

## 概述

goner/gin 的 HTTP 注入器（HttpInjector）是一个强大的依赖注入系统，能够自动将HTTP请求中的各种参数注入到处理函数的参数中。通过反射机制和标签解析，实现了类型安全的参数绑定和自动类型转换。

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

## 核心组件

### 1. HTTP 注入器 (HttpInjector)

- 负责将 HTTP 请求参数注入到处理函数参数中
- 支持从不同来源注入参数（Body、Header、Query、Param、Cookie）
- 通过反射实现自动参数绑定和类型转换

### 2. 延迟绑定注入器 (DelayBindInjector)

- 实现延迟绑定机制，在函数调用时才进行参数解析
- 提供高性能的参数注入能力
- 支持复杂的参数类型和嵌套结构

### 3. 绑定执行器 (BindExecutor)

- 管理类型解析器（TypeParser）和名称解析器（NameParser）
- 根据标签配置选择合适的解析器
- 执行具体的参数绑定逻辑

## 支持注入的类型和响应标签

### 类型解析器（TypeParser）支持的类型

这些类型不需要指定 `kind` 和 `key`，直接通过类型匹配进行注入：

| 名称        | 属性类型`${attributeType}` | 说明                                                           |
|-----------|------------------------|:-------------------------------------------------------------|
| **上下文注入** | `*gin.Context`         | 注入gin请求上下文指针。通过 `ginContextTypeParser` 实现。                   |
| **请求注入**  | `*http.Request`        | 注入http.Request指针。通过 `httpRequestTypeParser` 实现。              |
| **地址注入**  | `*url.URL`             | 注入url.URL指针。通过 `urlTypeParser` 实现。                           |
| **请求头注入** | `http.Header`          | 注入http.Header（请求头）。通过 `httpHeaderTypeParser` 实现。             |
| **响应注入**  | `gin.ResponseWriter`   | 注入gin.ResponseWriter（用于直接写入响应数据）。通过 `responseTypeParser` 实现。 |

### 名称解析器（NameParser）支持的类型

这些类型需要指定 `kind`，可选择性指定 `key`：

| 名称            | 属性类型`${attributeType}`                                        | 注入类型`${kind}` | 注入键值`${key}` | 说明                                                                                                                                              |
|---------------|---------------------------------------------------------------|:-------------:|:------------:|:------------------------------------------------------------------------------------------------------------------------------------------------|
| **Body注入**    | 结构体、结构体指针、[]byte、io.Reader、io.ReadCloser、any、Map、Slice、String |    `body`     |      /       | **body注入**；将请求body解析后注入到属性，注入类型为 `body`，不需要"注入键值`${key}`"；框架根据`Content-Type`自动判定是json还是xml等格式；每个请求处理函数只允许存在一个**body注入**。通过 `bodyNameParser` 实现。 |
| **请求头单值注入**   | number \| string                                              |    header     |    缺省取字段名    | 以键值`${key}`为`key`获取请求头，属性类型支持 简单类型<sub>[1]</sub>，解析不了会返回参数错误。通过 `headerNameParser` 实现。                                                          |
| **URL路径参数注入** | number \| string                                              |     param     |    缺省取字段名    | 以"注入键值`${key}`"为`key`调用函数`ctx.Param(key)`获取Url中定义的参数值，属性类型支持 简单类型<sub>[1]</sub>，解析不了会返回参数错误。通过 `paramNameParser` 实现。                            |
| **Query参数注入** | number \| string \| []number \| []string \| 结构体 \| 结构体指针      |     query     |    缺省取字段名    | 以"注入键值`${key}`"为`key`获取Query中的参数，属性类型支持 简单类型<sub>[1]</sub>，**支持简单类型的数组**，支持结构体和结构体指针，解析不了会返回参数错误。通过 `queryNameParser` 实现。                       |
| **Cookie注入**  | number \| string                                              |    cookie     |    缺省取字段名    | 以"注入键值`${key}`"为`key`调用函数`ctx.Cookie(key)`获取Cookie的值，属性类型支持 简单类型<sub>[1]</sub>，解析不了会返回参数错误。通过 `cookieNameParser` 实现。                             |

## 实现原理

### 工作流程

1. **加载阶段**：通过 `LoadGinHttpInjector` 函数加载HTTP注入器和相关解析器
2. **准备阶段**：`DelayBindInjector` 分析函数参数结构，为每个字段选择合适的解析器
3. **执行阶段**：当HTTP请求到达时，根据预编译的解析函数从请求中提取参数并注入

### 核心代码分析

```go
// LoadGinHttpInjector 加载HTTP注入器
func LoadGinHttpInjector(loader gone.Loader) error {
    loader.
        MustLoadX(injector.BuildLoad[*gin.Context](IdHttpInjector)).
        MustLoadX(parser.Load)
    return nil
}
```

这个函数完成了两个关键组件的加载：

1. `injector.BuildLoad[*gin.Context]` - 创建针对gin.Context的延迟绑定注入器
2. `parser.Load` - 加载所有的类型解析器和名称解析器

### 解析器架构

#### 类型解析器（TypeParser）

- `ginContextTypeParser` - 处理 `*gin.Context` 类型
- `httpRequestTypeParser` - 处理 `*http.Request` 类型
- `httpHeaderTypeParser` - 处理 `http.Header` 类型
- `urlTypeParser` - 处理 `*url.URL` 类型
- `responseTypeParser` - 处理 `gin.ResponseWriter` 类型

#### 名称解析器（NameParser）

- `bodyNameParser` - 处理 `body` 标签的参数注入
- `headerNameParser` - 处理 `header` 标签的参数注入
- `paramNameParser` - 处理 `param` 标签的参数注入
- `queryNameParser` - 处理 `query` 标签的参数注入
- `cookieNameParser` - 处理 `cookie` 标签的参数注入

## Query参数注入

### 属性类型为简单类型<sub>[1]</sub>

解析不了会返回参数错误。

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            yourName string `gone:"http,query=name"` //注册请求query中的name参数
            name string `gone:"http,query"`          //注册请求query中的name参数；不指定参数名，则取属性名作为参数名
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
        func (in struct {
            keywords []string `gone:"http,query=keyword"` //简单类型的query数组参数注入
        }) string {
            return fmt.Sprintf("hello, keywords is [%v]", in.keywords)
        },
    )
```

### 属性类型为结构体或者结构体指针

这种类型key无需指定；假设query为`?page=1&pageSize=20&keyword=gone&keyword=is&keyword=best`，`in.req`的值将会为
`{1,20,[]string{"gone","is","best"}}`；注意结构体中可以使用`form`标签进行属性映射。

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
        func (in struct {
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
        func (in struct {
            name string `gone:"http,param"`       //不指定参数名，使用属性名作为参数名
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

### 支持的Body类型

根据代码分析，`bodyNameParser` 支持以下类型：

1. **[]byte** - 直接读取原始字节数据
2. **io.Reader / io.ReadCloser** - 返回请求体的Reader接口
3. **结构体/Map/Slice/any** - 使用gin的ShouldBind进行自动绑定
4. **string** - 将请求体作为字符串读取

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
        func (in struct {
            req Req `gone:"http,body"` //注意：body只能被注入一次，因为 writer被读取后就变成空了
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
        func (in struct {
            contentType string `gone:"http,header"`              //不指定参数名，使用属性名作为参数名
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
        func (in struct {
            token string `gone:"http,cookie"`        //不指定参数名，使用属性名作为参数名
            token2 string `gone:"http,cookie=token"` //使用key指定参数名
        }) string {
            return "hello, your token in cookie is" + in.token
        },
    )
```

## 高级用法

### 类型解析器的直接注入

对于以下特殊类型，可以直接注入而无需指定 `kind` 和 `key`：

#### URL结构体注入

支持属性类型为 `*url.URL`，该类型定义在`net/url`包中，代表了HTTP请求的URL。

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            url *url.URL `gone:"http"` //使用结构体指针
        }) string {
            return "hello, your url is " + url.String()
        },
    )
```

#### 请求头完整注入

支持属性类型为 `http.Header`，该类型定义在`net/http`包中，代表了HTTP请求的所有Header。

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

#### 上下文结构体注入

支持属性类型为 `*gin.Context`，该类型定义在`github.com/gin-gonic/gin`包中，代表了HTTP请求的上下文。

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
            context *gin.Context `gone:"http"` //使用结构体指针
        }) string {
            return "hello, your method is " + context.Request.Method
        },
    )
```

#### 请求结构体注入

支持属性类型为 `*http.Request`，该类型定义在`net/http`包中，代表了HTTP请求信息。

```go
ctr.rootRouter.
    Group("/demo").
    POST(
        "/hello",
        func (in struct {
			request *http.Request `gone:"http"` //使用结构体指针
        }) string {
            return "hello, your method is " + request.Method
        },
    )
```

#### 请求响应接口注入

支持属性类型为 `gin.ResponseWriter`，该类型定义在`github.com/gin-gonic/gin`包中，代表了HTTP响应信息，可以使用该接口响应请求信息。

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

### 类型直接作为函数参数

goner/gin 的 HTTP 注入器不仅支持将参数包装在结构体中注入，还支持将特定类型直接作为函数参数使用。这种方式更加简洁，适用于只需要少量参数的场景。

支持直接作为函数参数的类型与类型解析器（TypeParser）支持的类型相同，包括：
- `*gin.Context` - HTTP 请求上下文
- `*http.Request` - HTTP 请求对象
- `*url.URL` - URL 对象
- `http.Header` - HTTP 请求头
- `gin.ResponseWriter` - HTTP 响应写入器

示例代码：

```go
ctr.rootRouter.
	Group("/demo").
	POST("/users", func(ctx *gin.Context, req *http.Request, writer gin.ResponseWriter){
		// 直接使用参数，无需从结构体中提取
		name := ctx.Query("name")
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte(fmt.Sprintf(`{"message":"Hello, %s"}`, name)))
    })
```

你也可以混合使用结构体参数和直接参数：

```go
ctr.rootRouter.
	Group("/demo").
	POST("/users", func(in struct {
		ID   int64  `gone:"http,param=id"`
		Name string `gone:"http,query=name"`
	}, writer gin.ResponseWriter){
		// 同时使用结构体参数和直接参数
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte(fmt.Sprintf(`{"id":%d, "name":"%s"}`, in.ID, in.Name)))
    })
```

这种方式使代码更加灵活，可以根据需要选择最合适的参数传递方式。


### 自定义参数解析器

goner/gin 允许开发者自定义参数解析器，以支持更复杂的参数注入场景。

#### 自定义类型解析器（TypeParser）

你可以通过实现 `injector.TypeParser[*gin.Context]` 接口来创建自定义的类型解析器。这个接口包含两个方法：

- `Parse(context *gin.Context) (reflect.Value, error)`: 从 `gin.Context` 中解析出目标类型的值。
- `Type() reflect.Type`: 返回解析器支持的目标类型。

**示例：自定义 Token 解析器**

假设你需要从请求头 `Authorization` 中解析出 `Bearer Token` 并注入到一个自定义的 `Token` 类型中。

1.  **定义 `Token` 类型：**

    ```go
    package main

    type Token string
    ```

2.  **实现 `TypeParser` 接口：**

    ```go
    package main

    import (
    	"reflect"
    	"strings"

    	"github.com/gin-gonic/gin"
    	"github.com/gone-io/gone/v2"
    	"github.com/gone-io/goner/gin/injector" // 确保导入正确的包
    )

    // 确保 tokenParser 实现了 TypeParser 接口
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
    	return reflect.Value{}, gone.NewParameterError("invalid token") // 使用 gone.NewParameterError 返回错误
    }

    func (t *tokenParser) Type() reflect.Type {
    	return reflect.TypeOf(Token(""))
    }
    ```

3.  **加载自定义解析器：**

    ```go
	gone.Load(&tokenParser{})
    ```

4.  **在 Handler 中使用：**

    现在你可以在你的 Gin Handler 中直接注入 `Token` 类型了。

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
    }) string { // 直接注入 Token 类型
    	return fmt.Sprintf("User token: %s", in.token)
    }

    func (ctr *MyController) Mount() gin.MountError {
    	ctr.Router.GET("/user", ctr.GetUser)
    	return nil
    }
    ```

    当请求 `/user` 接口并携带正确的 `Authorization: Bearer <your-token>` 请求头时，`token` 参数会被自动注入。

通过这种方式，你可以灵活地扩展 goner/gin 的参数注入能力，以适应各种复杂的业务需求。

#### 自定义名称解析器（NameParser）

除了自定义类型解析器外，goner/gin 还允许你创建自定义的名称解析器（NameParser）。名称解析器用于处理那些通过 `gone:"http,kind=key"` 标签指定的参数注入。这为你提供了更细粒度的控制，可以根据特定的 `kind` 和 `key` 来实现自定义的参数解析逻辑。

你需要实现 `injector.NameParser[*gin.Context]` 接口 来创建自定义的名称解析器。根据 `goner/gin/parser/name_query.go` 的示例，一个典型的名称解析器主要包含以下方法：

-   `Name() string`: 返回该名称解析器处理的 `kind` 类型。例如，`query`, `header`, `param`, `cookie`, `body`，或者你自定义的 `kind`。
-   `BuildParser(keyMap map[string]string, field reflect.StructField) (func(context *gin.Context) (reflect.Value, error), error)`: 这是核心的构建方法。它在应用初始化阶段被调用。
    -   `keyMap`: 一个从 `gone:"http,kind=key,..."` 标签中解析出来的键值对。例如，对于 `gone:"http,query=userId,optional"`，`keyMap` 可能包含 `{"query": "userId", "optional": ""}`。`s.Name()` (即 `kind`) 对应的值是主要的 `key`。
    -   `field`: 当前正在处理的结构体字段的 `reflect.StructField` 信息。
    -   该方法需要返回一个**解析函数**和一个错误。这个解析函数 `func(context *gin.Context) (reflect.Value, error)` 会在每个 HTTP 请求到达时被实际执行，用于从 `*gin.Context` 中提取数据、转换并返回 `reflect.Value`。

这种设计模式通过在初始化时预构建解析逻辑，优化了运行时性能。

**示例：自定义 CSV 解析器**

假设你需要从 query 参数中获取一个逗号分隔的字符串，并将其解析为一个字符串切片。例如，请求 URL 为 `/items?tags=go,gin,gone`，你希望将其注入到 `Tags []string` 字段中。

1.  **实现 `injector.NameParser[*gin.Context]`接口：**

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

    // 确保 csvQueryParser 实现了 NameParser 接口
    var _ injector.NameParser[*gin.Context] = (*csvQueryParser)(nil)

    type csvQueryParser struct {
    	gone.Flag
    }

    func (p *csvQueryParser) Name() string {
    	return "csv" // 自定义 kind 为 csv
    }

    func (p *csvQueryParser) BuildParser(keyMap map[string]string, field reflect.StructField) (func(context *gin.Context) (reflect.Value, error), error) {
    	// 从 keyMap 中获取 "csv" kind 对应的 key
    	// 例如，对于 gone:"http,csv=my_tags"，mainKey 应为 "my_tags"
    	// 对于 gone:"http,csv"，mainKey 可能为空，此时可以用字段名 field.Name
    	mainKey := keyMap[p.Name()] 
    	if mainKey == "" {
    		mainKey = field.Name // 如果tag中未指定key，则默认使用字段名
    	}

    	// 检查目标字段类型是否为 []string
    	if field.Type.Kind() != reflect.Slice || field.Type.Elem().Kind() != reflect.String {
    		return nil, fmt.Errorf("CSV parser: field '%s' must be of type []string, got %s", field.Name, field.Type.String())
    	}

    	// 返回实际执行解析的函数
    	return func(ctx *gin.Context) (reflect.Value, error) {
    		paramValue := ctx.Query(mainKey)
    		if paramValue == "" {
    			// 如果字段是必须的，这里可以返回 gone.NewParameterError
    			// 如果允许为空，则返回空切片
    			return reflect.ValueOf([]string{}), nil 
    		}
    		items := strings.Split(paramValue, ",")
    		return reflect.ValueOf(items), nil
    	}, nil
    }

    // 工厂函数用于被 gone.Load 加载
    func NewCsvQueryParser() gone.Goner {
    	return &csvQueryParser{}
    }
    ```


2.  **加载自定义解析器：**

    你需要将你的自定义名称解析器加载到 Gone 的依赖注入容器中。通常，这可以通过 `gone.Load()` 完成，并确保它被 `GinHttpInjector` 发现。
    `GinHttpInjector` 会收集所有实现了 `injector.NameParser[*gin.Context]` 接口的组件。

    ```go
    // 在你的 gone 应用启动逻辑中
    gone.Load(NewCsvQueryParser())
    ```

3.  **在 Handler 中使用：**

    现在你可以在你的 Gin Handler 的输入结构体中使用 `csv` 这个 `kind` 了。

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
    	Tags []string `gone:"http,csv=item_tags"` // 使用自定义的 csv 解析器
    	Name string   `gone:"http,query=name"`    // 使用内置的 query 解析器
    }

    func (ctr *MyController) CreateItem(in ItemRequest) string {
    	return fmt.Sprintf("Item created with name '%s' and tags: %v", in.Name, in.Tags)
    }

    func (ctr *MyController) Mount() error {
    	ctr.Router.POST("/items", ctr.CreateItem)
    	return nil
    }
    ```

    当请求 `/items?name=MyItem&item_tags=urgent,important` 时：
    - `in.Name` 会被注入为 `"MyItem"` (通过内置的 `queryNameParser`)
    - `in.Tags` 会被注入为 `[]string{"urgent", "important"}` (通过你的 `csvQueryParser`)

通过自定义名称解析器，你可以极大地增强 goner/gin 处理 HTTP 请求参数的灵活性和能力，使其适应各种复杂的 API 设计和数据格式。


## 性能优化

### 延迟绑定机制

goner/gin 的 HTTP 注入器采用延迟绑定机制，在应用启动时预编译参数解析函数，在请求处理时直接执行，避免了运行时的反射开销，提供了优秀的性能表现。

### 类型安全

所有的参数注入都是类型安全的，编译时就能发现类型不匹配的问题，运行时提供详细的错误信息帮助调试。

## 错误处理

当参数解析失败时，框架会返回 `gone.ParameterError`，包含详细的错误信息，帮助开发者快速定位问题。

## 备注

[1]. 简单类型指 字符串、布尔类型 和 数字类型，其中数字类型包括：

- 整数类型：int、uint、int8、uint8、int16、uint16、int32、uint32、int64、uint64
- 非负整数类型：uint、uint8、uint16、uint32、uint64
- 浮点类型：float32、float64
