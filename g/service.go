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
		Name:    name,
		Ip:      ip,
		Port:    port,
		Meta:    meta,
		Healthy: healthy,
		Weight:  weight,
	}
}

var _ Service = (*service)(nil)

type service struct {
	Name    string   `json:"name"`
	Ip      string   `json:"ip"`
	Port    int      `json:"port"`
	Meta    Metadata `json:"meta"`
	Healthy bool     `json:"healthy"`
	Weight  float64  `json:"weight"`
}

func (s *service) GetWeight() float64 {
	return s.Weight
}

func (s *service) GetName() string {
	return s.Name
}

func (s *service) GetIP() string {
	return s.Ip
}

func (s *service) GetPort() int {
	return s.Port
}

func (s *service) GetMetadata() Metadata {
	return s.Meta
}

func (s *service) IsHealthy() bool {
	return s.Healthy
}
