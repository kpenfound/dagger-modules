package main

import (
	"context"
	"fmt"
)

type Dockerd struct{}

func (t *Dockerd) Attach(
	ctx context.Context,
	c *Container,
	dockerVersion Optional[string],
) (*Container, error) {
	// docker service
	dockerV := dockerVersion.GetOr("24.0")
	port := 2375
	dockerd := dag.Container().
		From(fmt.Sprintf("docker:%s-dind", dockerV)).
		WithMountedCache(
			"/var/lib/docker",
			dag.CacheVolume(dockerV+"-docker-lib"),
			ContainerWithMountedCacheOpts{
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

	// Get dockerd endpoint
	dockerHost, err := dockerd.Endpoint(ctx, ServiceEndpointOpts{
		Scheme: "tcp",
	})
	if err != nil {
		return nil, err
	}

	return c.
		WithServiceBinding("docker", dockerd).
		WithEnvVariable("DOCKER_HOST", dockerHost), nil
}
