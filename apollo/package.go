package apollo

//go:generate sh -c "mockgen -package=apollo github.com/apolloconfig/agollo/v4 Client > agollo_mock.go"

//go:generate sh -c "mockgen -package=apollo github.com/apolloconfig/agollo/v4/agcache CacheInterface > cache_mock.go"
