// Proxy examples in Go
package main

import "dagger/example/internal/dagger"

type Example struct{}

// Example for WithService function
func (m *Example) Proxy_WithService(service *dagger.Service) *dagger.Service {
	return dag.Proxy().
		WithService(
			service,     // Dagger service to proxy
			"MyService", // Name of the service
			8080,        // Port for the proxy to listen on
			80,          // Port for the proxy to forward to
		).Service()
}

// Example for Service function
func (m *Example) Proxy_Service(serviceA *dagger.Service, serviceB *dagger.Service) *dagger.Service {
	return dag.Proxy().
		WithService(
			serviceA,   // Dagger service to proxy
			"ServiceA", // Name of the service
			8080,       // Port for the proxy to listen on
			80,         // Port for the proxy to forward to
		).
		WithService(
			serviceB,   // Dagger service to proxy
			"ServiceB", // Name of the service
			8081,       // Port for the proxy to listen on
			80,         // Port for the proxy to forward to
		).
		Service() // Return a Dagger service proxying to multiple services
}
