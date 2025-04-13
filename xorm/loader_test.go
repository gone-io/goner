package xorm

import (
	mock "github.com/gone-io/gone/mock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLoad(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	loader := mock.NewMockLoader(controller)
	loader.EXPECT().Loaded(gomock.Any()).Return(false)
	loader.EXPECT().Load(gomock.Any()).Return(nil)
	loader.EXPECT().Load(gomock.Any()).Return(nil)
	loader.EXPECT().Load(gomock.Any()).Return(nil)
	loader.EXPECT().Load(gomock.Any(), gomock.Any()).Return(nil)

	err := Load(loader)
	assert.Nil(t, err)
}
