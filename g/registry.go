package g

// ServiceRegistry provides methods for service registration and management
// It handles service instance registration, deregistration, and heartbeat updates
type ServiceRegistry interface {
	// Register adds a new service instance to the registry
	// Returns an error if registration fails
	Register(instance Service) error

	// Deregister removes a service instance from the registry
	// Returns an error if de registration fails
	Deregister(instance Service) error
}
