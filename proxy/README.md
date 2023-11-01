## Usage:

```go
return dag.Proxy().
	WithService("backend", 8080, 8080, backendService).
	WithService("frontend", 8081, 80, frontendService).
	Service()
```

