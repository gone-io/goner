package g

type ServiceRouter interface {
	GetServiceAddress(serviceName string) (serviceAddress string, err error)
}
