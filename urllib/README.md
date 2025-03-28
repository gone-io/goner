# Gone URLlib Component

`gone-urllib` is the HTTP client component of the Gone framework, implemented based on [imroc/req](https://github.com/imroc/req). It provides a simple and easy-to-use HTTP request functionality. With this component, you can easily make HTTP requests, handle responses, and integrate with other components in Gone applications.

## Features

- Seamless integration with Gone framework
- Clean API design
- Support for common HTTP methods (GET, POST, PUT, DELETE, etc.)
- Support for JSON request and response handling
- Support for request parameters, headers, and cookies
- Support for timeout control and retry mechanisms
- Support for HTTP/HTTPS proxy

## Installation

```bash
go get github.com/gone-io/goner/urllib
```

## Quick Start

### 1. Using URLlib in Your Application

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
	*req.Request `gone:"*"` // Inject HTTP client

	//*req.Client `gone:"*"` // Inject *req.Client

	//urllib.Client `gone:"*"` // Inject urllib.Client interface
}

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

### 2. Handling JSON Response

For more details, please refer to [imroc/req](https://github.com/imroc/req)

```go
func (s *MyService) GetUserInfo(userId string) (*UserInfo, error) {
    var result urllib.Res[UserInfo]  // Use generic response structure
    
    // Make request and parse JSON response
    resp, err := s.client.R().
        SetQueryParam("id", userId).
        SetResult(&result).  // Set response result
        Get("https://api.example.com/users")
    
    if err != nil {
        return nil, err
    }
    
    if !resp.IsSuccess() {
        return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
    }
    
    // Check business status code
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

### 3. Custom Client Configuration

```go
func (s *MyService) CustomizeClient() {
    // Get underlying req.Client for custom configuration
    client := s.client.C()
    
    // Set timeout
    client.SetTimeout(10 * time.Second)
    
    // Set retry
    client.SetCommonRetryCount(3)
    client.SetCommonRetryInterval(func(resp *req.Response, attempt int) time.Duration {
        return time.Duration(attempt) * time.Second
    })
    
    // Set proxy
    client.SetProxyURL("http://proxy.example.com:8080")
    
    // Set common headers
    client.SetCommonHeader("User-Agent", "Gone-URLlib/1.0")
}
```

## API Reference

### Client Interface

```go
type Client interface {
    // R creates a new request object
    R() *req.Request
    
    // C gets the underlying req.Client object
    C() *req.Client
}
```

### Res Structure

```go
type Res[T any] struct {
    Code int    `json:"code"`
    Msg  string `json:"msg,omitempty"`
    Data T      `json:"data,omitempty"`
}
```

A generic structure for handling standard JSON response format.

## Best Practices

1. Use dependency injection to get URLlib client, avoid manual creation
2. Create dedicated client wrappers for different API services to improve code reusability
3. Use generic response structure `Res<T>` to handle standard JSON response format
4. Set appropriate timeout and retry strategies to improve request reliability
5. Add trace ID in requests for better problem tracking

## Important Notes

1. When handling sensitive information (like authentication credentials), avoid hardcoding them in the code, use configuration or environment variables instead
2. For large file uploads or downloads, consider using streaming processing
3. In production environment, it's recommended to configure HTTPS certificate verification