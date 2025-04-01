package g

// Metadata represents a map of key-value pairs for storing service instance metadata
type Metadata map[string]string

// Service represents a service instance in the service registry
// It provides basic information about a service instance including its identity,
// location, metadata, and health status
type Service interface {
	// GetName returns the service name of the instance
	GetName() string

	GetIP() string

	GetPort() int

	// GetMetadata returns the metadata associated with the service instance
	GetMetadata() Metadata

	GetWeight() float64

	// IsHealthy returns the health status of the service instance
	IsHealthy() bool
}

func NewService(name, ip string, port int, meta Metadata, healthy bool, weight float64) Service {
	return &service{
		name:    name,
		ip:      ip,
		port:    port,
		meta:    meta,
		healthy: healthy,
		weight:  weight,
	}
}

var _ Service = (*service)(nil)

type service struct {
	name    string
	ip      string
	port    int
	meta    Metadata
	healthy bool
	weight  float64
}

func (s *service) GetWeight() float64 {
	return s.weight
}

func (s *service) GetName() string {
	return s.name
}

func (s *service) GetIP() string {
	return s.ip
}

func (s *service) GetPort() int {
	return s.port
}

func (s *service) GetMetadata() Metadata {
	return s.meta
}

func (s *service) IsHealthy() bool {
	return s.healthy
}
