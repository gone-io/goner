package xorm

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
	loader.EXPECT().MustLoad(gomock.Any()).Return(loader)
	loader.EXPECT().MustLoad(gomock.Any()).Return(loader)
	loader.EXPECT().MustLoad(gomock.Any()).Return(loader)
	loader.EXPECT().MustLoad(gomock.Any(), gomock.Any()).Return(loader)

	err := Load(loader)
	assert.Nil(t, err)
}
