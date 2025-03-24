package remote

import "github.com/spf13/viper"

type ViperInterface interface {
	ReadRemoteConfig() error
	AllSettings() map[string]any
	Get(key string) any
	MergeConfigMap(settings map[string]any) error
	UnmarshalKey(key string, rawVal any, opts ...viper.DecoderConfigOption) error
	SetConfigType(configType string)
	AddRemoteProvider(provider string, endpoint string, path string) error
	AddSecureRemoteProvider(provider string, endpoint string, path string, keyring string) error
}
