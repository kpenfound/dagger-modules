package main

import (
	"context"
	"fmt"
)

// Module for running docker in dagger
type Dockerd struct{}

// Attach a dockerd service to a container
func (t *Dockerd) Attach(
	ctx context.Context,
	container *Container,
	dockerVersion Optional[string],
) (*Container, error) {
	dockerd := t.Service(dockerVersion)

	dockerHost, err := dockerd.Endpoint(ctx, ServiceEndpointOpts{
		Scheme: "tcp",
	})
	if err != nil {
		return nil, err
	}

	return container.
		WithServiceBinding("docker", dockerd).
		WithEnvVariable("DOCKER_HOST", dockerHost), nil
}

// Get a Service container running dockerd
func (t *Dockerd) Service(dockerVersion Optional[string]) *Service {
	dockerV := dockerVersion.GetOr("24.0")
	port := 2375
	return dag.Container().
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
}
