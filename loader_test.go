package goner

import (
	mock "github.com/gone-io/gone/mock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGinLoad(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	loader := mock.NewMockLoader(controller)
	loader.EXPECT().Loaded(gomock.Any()).Return(false).AnyTimes()
	loader.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	loader.EXPECT().Load(gomock.Any(), gomock.Any()).AnyTimes()
	err := GinLoad(loader)
	assert.Nil(t, err)
}
