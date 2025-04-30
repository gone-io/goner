package gin

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func() {
			resp, err := http.Get("http://localhost:8080/metrics")
			assert.Nil(t, err)
			assert.Equalf(t, http.StatusOK, resp.StatusCode, "status code should be 200")
			all, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			assert.Contains(t, string(all), "go_gc_duration_seconds")
		})
}
