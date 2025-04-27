[//]: # (desc: http 客户端示例，使用goner/urllib)

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# http 客户端示例

这个示例展示了如何在Gone框架中使用URLlib组件发起HTTP请求。URLlib组件是基于[req](https://github.com/imroc/req)库封装的HTTP客户端组件。

## 快速开始

1. 首先，确保你的项目已经引入了必要的依赖：
```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/urllib"
    "github.com/imroc/req/v3"
)
```

2. 创建一个服务结构体，并注入HTTP客户端：
```go
type MyService struct {
    gone.Flag
    *req.Request `gone:"*"` // 注入 *req.Request
}
```

3. 实现你的业务方法：
```go
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
```

4. 在main函数中加载URLlib组件并运行：
```go
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

## 注入方式说明

URLlib组件提供了多种注入方式，你可以根据需要选择：

1. 注入 `*req.Request`（推荐）：
```go
*req.Request `gone:"*"` // 注入 *req.Request
```
这种方式直接注入请求对象，适合单次请求场景。

2. 注入 `*req.Client`：
```go
*req.Client `gone:"*"` // 注入 *req.Client
```
这种方式注入客户端对象，适合需要复用客户端配置的场景。

3. 注入 `urllib.Client` 接口：
```go
urllib.Client `gone:"*"` // 注入 urllib.Client 接口
```
这种方式注入客户端接口，适合需要mock测试的场景。

## 更多功能

- 支持设置请求头、查询参数、请求体等
- 支持文件上传、下载
- 支持自定义中间件
- 支持设置超时、重试等配置

更多使用方法请参考 [req](https://github.com/imroc/req) 文档。
