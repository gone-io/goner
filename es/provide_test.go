package es

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gone-io/gone/v2"
	"os"
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
