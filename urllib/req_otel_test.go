package urllib

import (
	mock "github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestWithOtel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	logger := mock.NewMockLogger(controller)

	r2 := r{
		logger:          logger,
		isOtelLogLoaded: true,
	}
	c := r2.C()
	_, err := c.R().Get("http://localhost:8080")
	assert.Error(t, err)
}
