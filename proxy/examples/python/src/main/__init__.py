
"""proxy examples in Python"""
import dagger
from dagger import dag, function, object_type, Service, Proxy

@object_type
class Example:

	@function
	def proxy_with_service(self, service: Service) -> Proxy:
		"""Example for with_service function"""
		return (
    		dag.proxy()
    		.with_service(
          		service,        # Dagger service to proxy
          		"my_service",   # Name of the service
                8080,           # Port for the proxy to listen on
                80              # Port for the proxy to forward to
    		)
		)


	@function
	def proxy_service(self, service_a: Service, service_b: Service) -> Service:
		"""Example for service function"""
		return (
    		dag.proxy()
    		.with_service(
          		service_a,      # Dagger service to proxy
          		"service_a",    # Name of the service
                8080,           # Port for the proxy to listen on
                80              # Port for the proxy to forward to
    		)
            .with_service(
          		service_b,      # Dagger service to proxy
          		"service_b",    # Name of the service
                8081,           # Port for the proxy to listen on
                80              # Port for the proxy to forward to
    		)
            .service() # Return a Dagger service proxying to multiple services
		)
