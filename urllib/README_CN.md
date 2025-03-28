# Gone URLlib 组件

`gone-urllib` 是 Gone 框架的 HTTP 客户端组件，基于 [imroc/req](https://github.com/imroc/req) 实现，提供了简洁易用的 HTTP 请求功能。通过该组件，您可以轻松地在 Gone 应用中发起 HTTP 请求，处理响应，并与其他组件集成。

## 功能特性

- 与 Gone 框架无缝集成
- 简洁的 API 设计
- 支持常见的 HTTP 方法（GET、POST、PUT、DELETE 等）
- 支持请求和响应的 JSON 处理
- 支持请求参数、头部和 Cookie 设置
- 支持超时控制和重试机制
- 支持 HTTP/HTTPS 代理

## 安装

```bash
go get github.com/gone-io/goner/urllib
```

## 快速开始

### 1. 在应用中使用 URLlib

```go
package main

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/urllib"
	"github.com/imroc/req/v3"
)

type MyService struct {
	gone.Flag
	*req.Request `gone:"*"` // 注入 HTTP 客户端

	//*req.Client `gone:"*"` // 注入 *req.Client

	//urllib.Client `gone:"*"` // 注入 urllib.Client 接口
}

func (s *MyService) GetData() (string, error) {
	// 发起 GET 请求
	resp, err := s.
		SetHeader("Accept", "application/json").
		Get("https://ipinfo.io")
	if err != nil {
		return "", err
	}

	// 获取响应内容
	return resp.String(), nil
}

func main() {
	gone.
		Load(&MyService{}).
		Loads(urllib.Load). // 加载 URLlib 组件
		Run(func(s *MyService) {
			data, err := s.GetData()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Data:", data)
		})
}

```

### 2. 处理 JSON 响应

更多说明，请参考 [imroc/req](https://github.com/imroc/req)

```go
func (s *MyService) GetUserInfo(userId string) (*UserInfo, error) {
    var result urllib.Res[UserInfo]  // 使用泛型响应结构
    
    // 发起请求并解析 JSON 响应
    resp, err := s.client.R().
        SetQueryParam("id", userId).
        SetResult(&result).  // 设置响应结果
        Get("https://api.example.com/users")
    
    if err != nil {
        return nil, err
    }
    
    if !resp.IsSuccess() {
        return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
    }
    
    // 检查业务状态码
    if result.Code != 0 {
        return nil, fmt.Errorf("business error: %s", result.Msg)
    }
    
    return &result.Data, nil
}

type UserInfo struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

### 3. 自定义客户端配置

```go
func (s *MyService) CustomizeClient() {
    // 获取底层的 req.Client 进行自定义配置
    client := s.client.C()
    
    // 设置超时
    client.SetTimeout(10 * time.Second)
    
    // 设置重试
    client.SetCommonRetryCount(3)
    client.SetCommonRetryInterval(func(resp *req.Response, attempt int) time.Duration {
        return time.Duration(attempt) * time.Second
    })
    
    // 设置代理
    client.SetProxyURL("http://proxy.example.com:8080")
    
    // 设置通用头部
    client.SetCommonHeader("User-Agent", "Gone-URLlib/1.0")
}
```

## API 参考

### Client 接口

```go
type Client interface {
    // R 创建一个新的请求对象
    R() *req.Request
    
    // C 获取底层的 req.Client 对象
    C() *req.Client
}
```

### Res 结构体

```go
type Res[T any] struct {
    Code int    `json:"code"`
    Msg  string `json:"msg,omitempty"`
    Data T      `json:"data,omitempty"`
}
```

用于处理标准 JSON 响应格式的泛型结构体。

## 最佳实践

1. 使用依赖注入获取 URLlib 客户端，避免手动创建
2. 为不同的 API 服务创建专门的客户端封装，提高代码复用性
3. 使用泛型响应结构 `Res<T>` 处理标准格式的 JSON 响应
4. 设置合理的超时和重试策略，提高请求可靠性
5. 在请求中添加追踪 ID，便于问题排查

## 注意事项

1. 处理敏感信息（如认证凭据）时，避免将其硬编码在代码中，推荐使用配置或环境变量
2. 对于大文件上传或下载，考虑使用流式处理方式
3. 在生产环境中，建议配置 HTTPS 证书验证