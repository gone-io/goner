package remote

import (
	"errors"
	"github.com/gone-io/gone/v2"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_configure_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLocalConfigure := NewMockConfigure(ctrl)
	mockViper := NewMockViperInterface(ctrl)

	tests := []struct {
		name                      string
		watch                     bool
		viperIsNil                bool
		viperGetValue             any
		unmarshalErr              error
		useLocalConfIfKeyNotExist bool
		expectLocalConfigureCall  bool
		expectErr                 error
	}{
		{
			name:                     "watch mode should record key",
			watch:                    true,
			viperGetValue:            "test value",
			unmarshalErr:             nil,
			expectLocalConfigureCall: false,
			expectErr:                nil,
		},
		{
			name:                     "viper is nil should fallback to local",
			viperIsNil:               true,
			expectLocalConfigureCall: true,
			expectErr:                nil,
		},
		{
			name:                     "empty viper value should fallback to local",
			viperGetValue:            "",
			expectLocalConfigureCall: true,
			expectErr:                nil,
		},
		{
			name:                     "normal viper value",
			viperGetValue:            "test value",
			unmarshalErr:             nil,
			expectLocalConfigureCall: false,
			expectErr:                nil,
		},
		{
			name:                      "unmarshal error with fallback enabled",
			viperGetValue:             "test value",
			unmarshalErr:              errors.New("unmarshal error"),
			useLocalConfIfKeyNotExist: true,
			expectLocalConfigureCall:  true,
			expectErr:                 nil,
		},
		{
			name:                      "unmarshal error with fallback disabled",
			viperGetValue:             "test value",
			unmarshalErr:              errors.New("unmarshal error"),
			useLocalConfIfKeyNotExist: false,
			expectLocalConfigureCall:  false,
			expectErr:                 gone.ToError(errors.New("unmarshal error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configure := &remoteConfigure{
				localConfigure:            mockLocalConfigure,
				keyMap:                    make(map[string][]any),
				watch:                     tt.watch,
				useLocalConfIfKeyNotExist: tt.useLocalConfIfKeyNotExist,
			}

			if !tt.viperIsNil {
				configure.viper = mockViper
				mockViper.EXPECT().Get("test.key").Return(tt.viperGetValue)
				if tt.viperGetValue != "" {
					mockViper.EXPECT().UnmarshalKey("test.key", gomock.Any()).Return(tt.unmarshalErr)
				}
			}

			if tt.expectLocalConfigureCall {
				mockLocalConfigure.EXPECT().Get("test.key", gomock.Any(), "default").Return(nil)
			}

			var value string
			err := configure.Get("test.key", &value, "default")

			if tt.watch {
				if _, exists := configure.keyMap["test.key"]; !exists {
					t.Error("key should be recorded in watch mode")
				}
			}

			if (err == nil && tt.expectErr != nil) || (err != nil && tt.expectErr == nil) ||
				(err != nil && tt.expectErr != nil && err.(gone.Error).Code() != tt.expectErr.(gone.Error).Code()) {
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
			}
		})
	}
}
