package gin

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Load(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	t.Run("Success", func(t *testing.T) {
		mockLoader := NewMockLoader(controller)
		mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)

		// 期望加载 router
		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		// 期望加载 SysMiddleware
		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		// 期望加载 proxy
		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		// 期望加载 GinResponser
		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		// 期望加载 httpInjector
		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		// 期望加载 GinServer
		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		err := Load(mockLoader)
		assert.NoError(t, err)
	})

	t.Run("Error_LoadRouter", func(t *testing.T) {
		mockLoader := NewMockLoader(controller)
		mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)
		expectedErr := errors.New("load router error")

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		err := Load(mockLoader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load router error")
	})

	t.Run("Error_LoadSysMiddleware", func(t *testing.T) {
		mockLoader := NewMockLoader(controller)
		mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)
		expectedErr := errors.New("load sys middleware error")

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(expectedErr)

		err := Load(mockLoader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load sys middleware error")
	})

	t.Run("Error_LoadProxy", func(t *testing.T) {
		mockLoader := NewMockLoader(controller)
		mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)
		expectedErr := errors.New("load proxy error")

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		err := Load(mockLoader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load proxy error")
	})

	t.Run("Error_LoadGinResponser", func(t *testing.T) {
		mockLoader := NewMockLoader(controller)
		mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)
		expectedErr := errors.New("load gin responser error")

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(expectedErr)

		err := Load(mockLoader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load gin responser error")
	})

	t.Run("Error_LoadHttpInjector", func(t *testing.T) {
		mockLoader := NewMockLoader(controller)
		mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)
		expectedErr := errors.New("load http injector error")

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(expectedErr)

		err := Load(mockLoader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load http injector error")
	})

	t.Run("Error_LoadGinServer", func(t *testing.T) {
		mockLoader := NewMockLoader(controller)
		mockLoader.EXPECT().Loaded(gomock.Any()).Return(false)
		expectedErr := errors.New("load gin server error")

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any()).
			Return(nil)

		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		err := Load(mockLoader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load gin server error")
	})
}
