package remote

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_remoteConfigure_Notify(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockViper := NewMockViperInterface(controller)

	newRemoteViper = func() ViperInterface {
		return mockViper
	}
	mockViper.EXPECT().SetConfigType(gomock.Any())
	mockViper.EXPECT().AddRemoteProvider(gomock.Any(), gomock.Any(), gomock.Any())
	mockViper.EXPECT().ReadRemoteConfig()
	mockViper.EXPECT().AllSettings().Return(map[string]any{
		"test": "test",
	})

	gone.
		NewApp(Load).
		Test(func(watcher gone.ConfWatcher, w *watcher) {
			key := "test"
			var oldVal, newVal any
			watcher(key, func(o, n any) {
				oldVal, newVal = o, n
			})

			mockViper.EXPECT().ReadRemoteConfig()
			mockViper.EXPECT().AllSettings().Return(map[string]any{
				"test": "test2",
			})
			w.doWatch()
			assert.Equal(t, "test", oldVal)
			assert.Equal(t, "test2", newVal)

		})
}
