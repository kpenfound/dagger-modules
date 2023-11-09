package main

import (
	"context"
	"fmt"
)

type Testcontainers struct{}

// optionally, boolean for Testcontainers Cloud, optional TCC_TOKEN (secret).
func (t *Testcontainers) Enable(ctx context.Context, c *Container) (*Container, error) {
	// docker service
	dockerVersion := "24.0" // extract as a parameter
	port := 2375
	dockerd := dag.Container().From(fmt.Sprintf("docker:%s-dind", dockerVersion)).
		WithMountedCache("/var/lib/docker", dag.CacheVolume(dockerVersion+"-docker-lib"), ContainerWithMountedCacheOpts{
			Sharing: Private,
		}).
		WithExposedPort(port).
		WithExec([]string{
			"dockerd",
			"--host=tcp://0.0.0.0:2375",
			"--host=unix:///var/run/docker.sock",
			"--tls=false",
		}, ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		AsService()

	dockerHost, err := dockerd.Endpoint(ctx, ServiceEndpointOpts{
		Scheme: "tcp",
	})
	if err != nil {
		return nil, err
	}
	// ---------

	return c.
		WithServiceBinding("docker", dockerd).
		WithEnvVariable("DOCKER_HOST", dockerHost), nil

}
