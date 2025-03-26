// Package es provides integration with Elasticsearch for gone applications.
// It offers functionality to create and manage Elasticsearch client instances
// in a gone-compatible way, supporting both regular and typed clients.
package es

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gone-io/gone/v2"
)

// Load registers a singleton Elasticsearch client provider with the gone loader.
// It ensures only one client instance is created and reused across the application.
//
// Parameters:
//   - loader: The gone.Loader instance used for dependency injection
//
// Returns:
//   - error: Any error encountered during client creation or registration
func Load(loader gone.Loader) error {
	var load = gone.OnceLoad(func(loader gone.Loader) error {
		var single *elasticsearch.Client

		getSingleEs := func(
			tagConf string,
			param struct {
				config elasticsearch.Config `gone:"config,es"`
			},
		) (*elasticsearch.Client, error) {
			var err error
			if single == nil {
				if single, err = elasticsearch.NewClient(param.config); err != nil {
					return nil, gone.ToError(err)
				}
			}
			return single, nil
		}
		provider := gone.WrapFunctionProvider(getSingleEs)
		return loader.Load(provider)
	})
	return load(loader)
}

// LoadTypedClient registers a singleton TypedClient provider with the gone loader.
// TypedClient provides a more type-safe way to interact with Elasticsearch.
//
// Parameters:
//   - loader: The gone.Loader instance used for dependency injection
//
// Returns:
//   - error: Any error encountered during client creation or registration
func LoadTypedClient(loader gone.Loader) error {
	var load = gone.OnceLoad(func(loader gone.Loader) error {
		var single *elasticsearch.TypedClient

		getSingleEs := func(
			tagConf string,
			param struct {
				config elasticsearch.Config `gone:"config,es"`
			},
		) (*elasticsearch.TypedClient, error) {
			var err error
			if single == nil {
				if single, err = elasticsearch.NewTypedClient(param.config); err != nil {
					return nil, gone.ToError(err)
				}
			}
			return single, nil
		}
		provider := gone.WrapFunctionProvider(getSingleEs)
		return loader.Load(provider)
	})
	return load(loader)
}
