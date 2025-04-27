[//]: # (desc: http client example using goner/urllib)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# HTTP Client Example

This example demonstrates how to make HTTP requests using the URLlib component in the Gone framework. The URLlib component is an HTTP client component based on the [req](https://github.com/imroc/req) library.

## Quick Start

1. First, ensure that your project has imported the necessary dependencies:
```go
import (
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/urllib"
    "github.com/imroc/req/v3"
)
```

2. Create a service struct and inject the HTTP client:
```go
type MyService struct {
    gone.Flag
    *req.Request `gone:"*"` // Inject *req.Request
}
```

3. Implement your business method:
```go
func (s *MyService) GetData() (string, error) {
    // Make a GET request
    resp, err := s.
        SetHeader("Accept", "application/json").
        Get("https://ipinfo.io")
    if err != nil {
        return "", err
    }

    // Get response content
    return resp.String(), nil
}
```

4. Load the URLlib component and run in the main function:
```go
func main() {
    gone.
        Load(&MyService{}).
        Loads(urllib.Load). // Load URLlib component
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

## Injection Methods

The URLlib component provides multiple injection methods that you can choose based on your needs:

1. Inject `*req.Request` (Recommended):
```go
*req.Request `gone:"*"` // Inject *req.Request
```
This method directly injects the request object, suitable for single request scenarios.

2. Inject `*req.Client`:
```go
*req.Client `gone:"*"` // Inject *req.Client
```
This method injects the client object, suitable for scenarios where client configuration needs to be reused.

3. Inject `urllib.Client` interface:
```go
urllib.Client `gone:"*"` // Inject urllib.Client interface
```
This method injects the client interface, suitable for mock testing scenarios.

## More Features

- Support for setting request headers, query parameters, request body, etc.
- Support for file upload and download
- Support for custom middleware
- Support for timeout and retry configuration

For more usage methods, please refer to the [req](https://github.com/imroc/req) documentation.