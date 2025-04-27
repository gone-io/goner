package gin

import (
	"github.com/gone-io/gone/mock/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Load(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	t.Run("Success", func(t *testing.T) {
		mockLoader := mock.NewMockLoader(controller)

		// 期望加载 router
		mockLoader.EXPECT().
			MustLoad(gomock.Any(), gomock.Any()).
			Return(mockLoader)

		// 期望加载 SysMiddleware
		mockLoader.EXPECT().
			MustLoad(gomock.Any()).
			Return(mockLoader)

		// 期望加载 proxy
		mockLoader.EXPECT().
			MustLoad(gomock.Any(), gomock.Any()).
			Return(mockLoader)

		// 期望加载 GinResponser
		mockLoader.EXPECT().
			MustLoad(gomock.Any()).
			Return(mockLoader)

		// 期望加载 httpInjector
		mockLoader.EXPECT().
			MustLoad(gomock.Any()).
			Return(mockLoader)

		// 期望加载 GinServer
		mockLoader.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil)

		err := Load(mockLoader)
		assert.NoError(t, err)
	})
}
