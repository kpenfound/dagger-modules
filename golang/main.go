// Build Go projects
//
// A utility module for building, testing, and linting Go projects

package main

import (
	"context"
	"fmt"
	"runtime"

	"golang/internal/dagger"
)

const (
	MOUNT_PATH = "/app"
	LINT_IMAGE = "golangci/golangci-lint:v2.1"
	OUT_DIR    = "/out/"
)

type Golang struct {
	// Golang container with the go project and go toolchain
	Container *dagger.Container
}

func New(
	// Specify an alternative container to index.docker.io/golang
	// +optional
	ctr *dagger.Container,
	// The source directory of the Go project
	// +optional
	source *dagger.Directory,
	// Golang image tag to use
	// +default="1.24"
	version string,
) *Golang {
	g := &Golang{}
	if ctr == nil {
		ctr = g.Base(version).Container
	}
	g.Container = ctr

	if source != nil {
		g.Container = g.Container.WithDirectory(MOUNT_PATH, source)
	}

	return g
}

// Build the Go project
func (g *Golang) Build(
	// The Go source code to build
	// +optional
	source *dagger.Directory,
	// Arguments to `go build`
	args []string,
	// The architecture for GOARCH
	// +optional
	arch string,
	// The operating system for GOOS
	// +optional
	os string,
) *dagger.Directory {
	if arch == "" {
		arch = runtime.GOARCH
	}
	if os == "" {
		os = runtime.GOOS
	}

	if source != nil {
		g = g.WithSource(source)
	}

	command := append([]string{"go", "build", "-o", OUT_DIR}, args...)
	return g.Container.
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("GOOS", os).
		WithExec(command).
		Directory(OUT_DIR)
}

// Build a Go project returning a Container containing the build
func (g *Golang) BuildContainer(
	// The Go source code to build
	// +optional
	source *dagger.Directory,
	// Arguments to `go build`
	// +optional
	args []string,
	// The architecture for GOARCH
	// +optional
	arch string,
	// The operating system for GOOS
	// +optional
	os string,
	// Base container in which to copy the build
	// +optional
	base *dagger.Container,
) *dagger.Container {
	dir := g.Build(source, args, arch, os)
	if base == nil {
		base = dag.Container().From("ubuntu:latest")
	}
	return base.
		WithDirectory("/usr/local/bin/", dir)
}

// Test the Go project
func (g *Golang) Test(
	ctx context.Context,
	// The Go source code to test
	// +optional
	source *dagger.Directory,
	// Arguments to `go test`
	// +optional
	// +default "./..."
	args []string,
) (string, error) {
	if source != nil {
		g = g.WithSource(source)
	}
	command := append([]string{"go", "test"}, args...)
	return g.Container.WithExec(command).Stdout(ctx)
}

// Format the Go project
func (g *Golang) Fmt(
	ctx context.Context,
	// The Go source code to format
	// +optional
	source *dagger.Directory,
	// Arguments to `go fmt`
	// +optional
	// +default "./..."
	args []string,
) (*Golang, error) {
	if source != nil {
		g = g.WithSource(source)
	}
	command := append([]string{"go", "fmt"}, args...)
	c, err := g.Container.WithExec(command).Sync(ctx)
	if err != nil {
		return nil, err
	}
	g.Container = c
	return g, err
}

// Lint the Go project
func (g *Golang) GolangciLint(
	ctx context.Context,
	// The Go source code to lint
	// +optional
	source *dagger.Directory,
) (string, error) {
	return g.lint(source, []string{}).
		Stdout(ctx)
}

// Lint the Go project and apply fixes
func (g *Golang) GolangciLintFix(
	ctx context.Context,
	// The Go source code to lint
	// +optional
	source *dagger.Directory,
) (*dagger.Directory, error) {
	done, err := g.lint(source, []string{"--fix"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	return done.Directory(MOUNT_PATH), nil
}

// Private lint helper
func (g *Golang) lint(source *dagger.Directory, args []string) *dagger.Container {
	if source == nil {
		source = g.GetSource()
	}
	return dag.Container().From(LINT_IMAGE).
		WithMountedDirectory(MOUNT_PATH, source).
		WithWorkdir(MOUNT_PATH).
		WithExec(append([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}, args...))
}

// Sets up the Container with a golang image and cache volumes
func (g *Golang) Base(
	// Golang image tag to use
	// +default="1.24"
	version string,
) *Golang {
	image := fmt.Sprintf("golang:%s", version)
	g.Container = dag.Container().
		From(image).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gomodcache")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("gobuildcache")).
		WithWorkdir(MOUNT_PATH)
	return g
}

// Specify the Source to use in the module
func (g *Golang) WithSource(source *dagger.Directory) *Golang {
	g.Container = g.Container.WithDirectory(MOUNT_PATH, source)
	return g
}

// Get the current state of the source directory
func (g *Golang) GetSource() *dagger.Directory {
	return g.Container.Directory(MOUNT_PATH)
}

// Bring your own container
func (g *Golang) WithContainer(ctr *dagger.Container) *Golang {
	g.Container = ctr
	return g
}
