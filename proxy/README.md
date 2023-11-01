## Usage:

```go
dag.Proxy().
	WithService(backendService, "backend", 8080, 8080).
	WithService(frontendService, "frontend", 8081, 80).
	Service()
```

