package remote

//go:generate mockgen -source=interface.go -destination=viper_interface_mock.go -package=remote
//go:generate mockgen -destination=configure_mock.go -package=remote github.com/gone-io/gone/v2 Configure,Logger
