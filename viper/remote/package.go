package remote

//go:generate mockgen -source=remote.go -destination=mock_viper_interface_test.go -package=remote
//go:generate mockgen -destination=mock_configure_test.go -package=remote github.com/gone-io/gone/v2 Configure,Logger
