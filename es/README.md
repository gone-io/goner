# Gone Elasticsearch Integration

This package provides Elasticsearch integration for Gone applications, offering both low-level and typed client support.

## Features

- Easy integration with Gone's dependency injection system
- Support for both low-level and typed Elasticsearch clients
- Singleton client instance management
- Comprehensive configuration options

## Installation

```bash
go get github.com/gone-io/goner/es
```

## Configuration

Create a `default.yaml` file in your project's config directory with the following Elasticsearch configuration:

```yaml
es:
  addresses: http://localhost:9200   # A list of Elasticsearch nodes to use
  username:   # Username for HTTP Basic Authentication
  password:   # Password for HTTP Basic Authentication

  cloudID:    # Endpoint for the Elastic Service
  aPIKey:     # Base64-encoded token for authorization
  serviceToken: # Service token for authorization

  # Additional optional configurations
  certificateFingerprint:  # SHA256 hex fingerprint
  retryOnStatus:          # List of status codes for retry (default: 502, 503, 504)
  maxRetries:             # Default: 3
  compressRequestBody:    # Default: false
  enableMetrics:          # Enable metrics collection
  enableDebugLogger:      # Enable debug logging
```

## Usage

### Low-Level Client

```go
package main

import (
    "bytes"
    "encoding/json"
    "github.com/elastic/go-elasticsearch/v8"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/es"
    "github.com/gone-io/goner/viper"
    "io"
)

type esUser struct {
    gone.Flag
    esClient *elasticsearch.Client `gone:"*"`
    logger   gone.Logger          `gone:"*"`
}

func (s *esUser) Use() {
    // Create index
    create, err := s.esClient.Indices.Create("my_index")
    if err != nil {
        s.logger.Errorf("Indices.Create err:%v", err)
        return
    }

    // Create document
    document := struct {
        Name string `json:"name"`
    }{
        "go-elasticsearch",
    }
    data, _ := json.Marshal(document)
    index, err := s.esClient.Index("my_index", bytes.NewReader(data))
    if err != nil {
        s.logger.Errorf("Index err:%v", err)
        return
    }

    // Get document ID
    var id struct {
        ID string `json:"_id"`
    }
    all, _ := io.ReadAll(index.Body)
    json.Unmarshal(all, &id)

    // Get document
    get, err := s.esClient.Get("my_index", id.ID)
    if err != nil {
        s.logger.Errorf("Get err:%v", err)
    }
}

func main() {
    gone.NewApp(
        viper.Load, // Load configuration
        es.Load,    // Initialize Elasticsearch client
    ).Run(func(esUser *esUser) {
        esUser.Use()
    })
}
```

### Typed Client

```go
package main

import (
    "context"
    "github.com/elastic/go-elasticsearch/v8"
    "github.com/gone-io/gone/v2"
    "github.com/gone-io/goner/es"
    "github.com/gone-io/goner/viper"
)

type esUser struct {
    gone.Flag
    esClient *elasticsearch.TypedClient `gone:"*"`
    logger   gone.Logger               `gone:"*"`
}

func (s *esUser) Use() {
    ctx := context.TODO()

    // Create index
    create, err := s.esClient.Indices.Create("my_index").Do(ctx)
    if err != nil {
        s.logger.Errorf("Indices.Create err:%v", err)
        return
    }

    // Create document
    document := struct {
        Name string `json:"name"`
    }{
        "go-elasticsearch",
    }
    index, err := s.esClient.Index("my_index").Document(document).Do(ctx)
    if err != nil {
        s.logger.Errorf("Index err:%v", err)
        return
    }

    // Get document
    get, err := s.esClient.Get("my_index", index.Id_).Do(ctx)
    if err != nil {
        s.logger.Errorf("Get err:%v", err)
    }
}

func main() {
    gone.NewApp(
        viper.Load, // Load configuration
        es.Load,    // Initialize Elasticsearch client
    ).Run(func(esUser *esUser) {
        esUser.Use()
    })
}
```

## API Reference

### Load

```go
func Load(loader gone.Loader) error
```

Registers a singleton low-level Elasticsearch client provider with the Gone loader.

### LoadTypedClient

```go
func LoadTypedClient(loader gone.Loader) error
```

Registers a singleton typed Elasticsearch client provider with the Gone loader.

## License

This project is licensed under the MIT License - see the LICENSE file for details.