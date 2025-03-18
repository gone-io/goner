package apollo

//go:generate sh -c "mockgen -package=apollo github.com/apolloconfig/agollo/v4 Client > agollo_mock_test.go"

//go:generate sh -c "mockgen -package=apollo github.com/gone-io/gone/v2 Configure,Logger > gone_mock_test.go"

//go:generate sh -c "mockgen -package=apollo github.com/apolloconfig/agollo/v4/agcache CacheInterface > cache_mock_test.go"
