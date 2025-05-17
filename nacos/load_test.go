package nacos

import (
	mock "github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLoad(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	loader := mock.NewMockLoader(controller)
	loader.EXPECT().Loaded(gomock.Any()).Return(false).AnyTimes()
	loader.EXPECT().Load(gomock.Any(), gomock.Any()).AnyTimes()

	err := Load(loader)
	assert.Nil(t, err)

	err = RegistryLoad(loader)
	assert.Nil(t, err)
}
