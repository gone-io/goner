package viper

//go:generate mockgen -destination=mock_watcher_keeper_test.go -package=viper . WatcherKeeper
//go:generate mockgen -destination=mock_key_getter_test.go -package=viper . KeyGetter
//go:generate mockgen -destination=mock_configure_test.go -package=remote github.com/gone-io/gone/v2 Configure
