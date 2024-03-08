// Utility for running dockerd in Dagger
//
// A utility module for configuring a dockerd service in your Dagger pipeline

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
	// +optional
	// +default="24.0"
	dockerVersion string,
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
func (t *Dockerd) Service(
	// +optional
	// +default="24.0"
	dockerVersion string,
) *Service {
	port := 2375
	return dag.Container().
		From(fmt.Sprintf("docker:%s-dind", dockerVersion)).
		WithMountedCache(
			"/var/lib/docker",
			dag.CacheVolume(dockerVersion+"-docker-lib"),
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
