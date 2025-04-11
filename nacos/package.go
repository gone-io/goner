package nacos

//go:generate mockgen -destination=nacos_mock.go -package=nacos github.com/nacos-group/nacos-sdk-go/v2/clients/config_client IConfigClient,
//go:generate mockgen -destination=nacos_name_mock.go -package=nacos github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client INamingClient
