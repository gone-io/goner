package g

import "context"

// LoadBalancer provides load balancing functionality for service instances
// It selects an appropriate instance from available service instances
type LoadBalancer interface {
	// GetInstance returns a service instance based on the load balancing strategy
	// Returns an error if no instance is available or selection fails
	GetInstance(ctx context.Context, serviceName string) (Service, error)
}

// LoadBalanceStrategy defines the interface for implementing different load balancing algorithms
// It allows for custom instance selection logic based on various factors
type LoadBalanceStrategy interface {
	// Select chooses a service instance from the provided list of instances
	// Returns an error if selection fails or no suitable instance is found
	Select(ctx context.Context, instances []Service) (Service, error)
}
