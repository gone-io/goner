package svc

// Endpoint interface for service.
type Endpoint interface {
	// Host returns the IPv4/IPv6 address of a service.
	Host() string

	// Port returns the port of a service.
	Port() int

	// String formats and returns the Endpoint as a string.
	String() string
}

// Endpoints are composed by multiple Endpoint.
type Endpoints []Endpoint

// Metadata stores custom key-value pairs.
type Metadata map[string]interface{}

type Service interface {
	// GetName returns the name of the service.
	// The name is necessary for a service, and should be unique among services.
	GetName() string

	// GetMetadata returns the Metadata map of service.
	// The Metadata is key-value pair map specifying extra attributes of a service.
	GetMetadata() Metadata

	// GetEndpoints returns the Endpoints of service.
	// The Endpoints contain multiple host/port information of service.
	GetEndpoints() Endpoints
}

// SearchInput is the input for service searching.
type SearchInput struct {
	Name     string   // Search by service name.
	Metadata Metadata // Filter by metadata if there are multiple result.
}

// Registrar interface for service registrar.
type Registrar interface {
	// Register registers `service` to Registry.
	// Note that it returns a new Service if it changes the input Service with custom one.
	Register(service Service) (registered Service, err error)

	// Deregister off-lines and removes `service` from the Registry.
	Deregister(service Service) error
}

// Discovery interface for service discovery.
type Discovery interface {
	// Search searches and returns services with specified condition.
	Search(in SearchInput) (result []Service, err error)

	// Watch watches specified condition changes.
	Watch(serviceName string) (watcher Watcher, err error)
}

// Watcher interface for service.
type Watcher interface {
	// Proceed proceeds watch in blocking way.
	// It returns all complete services that watched by `key` if any change.
	Proceed() (services []Service, err error)

	// Close closes the watcher.
	Close() error
}

// Registry interface for service.
type Registry interface {
	Registrar
	Discovery
}
