package g

// ServiceDiscovery provides methods for discovering and monitoring service instances
// It allows clients to find available service instances and watch for changes
type ServiceDiscovery interface {
	// GetInstances returns all instances of a specified service
	// Returns an error if the service discovery fails
	GetInstances(serviceName string) ([]Service, error)

	// Watch creates a channel that receives updates when the service instances change
	// Returns an error if watching fails
	Watch(serviceName string) (ch <-chan []Service, stop func() error, err error)
}
