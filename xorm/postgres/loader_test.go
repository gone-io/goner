package postgres

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLoad(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	loader := gone.NewMockLoader(controller)
	loader.EXPECT().MustLoadX(gomock.Any()).Return(loader)

	err := Load(loader)
	assert.Nil(t, err)
}
