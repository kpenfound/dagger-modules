// A generated module for Kind functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"fmt"
	"runtime"
)

type Kind struct{}

const (
	default_version = "v0.22.0"
	docker_host = "dockerd"
	docker_version = "24.0"
)

func (k *Kind) Cli(
	// +optional
	version string,
) *Container {
	if version == "" {
		version = default_version
	}

	binary := dag.HTTP(fmt.Sprintf("https://kind.sigs.k8s.io/dl/%s/kind-linux-%s", version, runtime.GOARCH))

	return dag.Container().From(fmt.Sprintf("index.docker.io/docker:%s-cli", docker_version)).
		WithFile("/bin/kind", binary).
		WithExec([]string{"chmod", "+x", "/bin/kind"})
}

func (k *Kind) Demo() *Container {
	config := `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "0.0.0.0"
kubeadmConfigPatches:
- |
  kind: ClusterConfiguration
  apiServer:
    certSANs:
      - "dockerd"
`
	ctr := k.Cli(default_version).
		WithNewFile("/kind_config.yml", ContainerWithNewFileOpts{ Contents: config }).
		WithEnvVariable("KIND_EXPERIMENTAL_PROVIDER", "docker").
		WithEnvVariable("DOCKER_HOST", fmt.Sprintf("tcp://%s:2375", docker_host)).
		WithServiceBinding(docker_host, dockerd()).
		WithExec([]string{"kind", "create", "cluster", "--retain", "--config", "/kind_config.yml"})

	return ctr
}

func dockerd() *Service {
	return dag.Container().
		From(fmt.Sprintf("index.docker.io/docker:%s-dind", docker_version)).
		WithoutEntrypoint().
		WithExposedPort(2375).
		WithExec([]string{
			"dockerd",
			"--host=tcp://0.0.0.0:2375",
			"--tls=false",
		}, ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		AsService()
}
