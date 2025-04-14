package es

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gone-io/gone/v2"
	"os"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	_ = os.Setenv("GONE_ES", `{"ADDRESSES":["http://127.0.0.1:9200"]}`)

	gone.
		NewApp(Load).
		Run(func(client *elasticsearch.Client) {
			_, err := client.Ping()
			if err != nil {
				t.Error(err)
			}
		})
}

func TestLoadTypedClient(t *testing.T) {
	_ = os.Setenv("GONE_ES", `{"ADDRESSES":["http://127.0.0.1:9200"]}`)

	gone.
		NewApp(LoadTypedClient).
		Run(func(client *elasticsearch.TypedClient) {
			_, err := client.Ping().Do(context.Background())
			if err != nil {
				t.Error(err)
			}
		})
}

func TestLoadError(t *testing.T) {

	t.Run("Load error", func(t *testing.T) {
		old := newClient
		defer func() {
			newClient = old
		}()
		newClient = func(config elasticsearch.Config) (*elasticsearch.Client, error) {
			return nil, errors.New("create-error")
		}

		defer func() {
			if err := recover(); err != nil {
				if !strings.Contains(err.(gone.Error).Error(), "create-error") {
					t.Error(err)
				}
			}
		}()
		gone.NewApp(Load).Run(func(*elasticsearch.Client) {})
	})

	t.Run("LoadTypedClient", func(t *testing.T) {
		old := newTypedClient
		defer func() {
			newTypedClient = old
		}()
		newTypedClient = func(cfg elasticsearch.Config) (*elasticsearch.TypedClient, error) {
			return nil, errors.New("create-error")
		}
		defer func() {
			if err := recover(); err != nil {
				if !strings.Contains(err.(gone.Error).Error(), "create-error") {
					t.Error(err)
				}
			}
		}()
		gone.NewApp(LoadTypedClient).Run(func(*elasticsearch.TypedClient) {})
	})
}
