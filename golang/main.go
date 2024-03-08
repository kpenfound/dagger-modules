// Build Go projects
//
// A utility module for building, testing, and linting Go projects

package main

import (
	"context"
	"fmt"
	"runtime"
)

const (
	DEFAULT_GO = "1.21"
	PROJ_MOUNT = "/src"
	LINT_IMAGE = "golangci/golangci-lint:v1.55.2"
)

type Golang struct {
	// +private
	Ctr *Container
	// +private
	Proj *Directory
}

func New(
	// +optional
	ctr *Container,
	// +optional
	proj *Directory,
) *Golang {
	g := &Golang{}
	if ctr == nil {
		ctr = g.Base(DEFAULT_GO).Ctr
	}
	g.Ctr = ctr

	if proj != nil {
		g.Proj = proj
	}

	return g
}

// Build the Go project
func (g *Golang) Build(
	// The Go source code to build
	// +optional
	source *Directory,
	// Arguments to `go build`
	args []string,
	// The architecture for GOARCH
	// +optional
	arch string,
	// The operating system for GOOS
	// +optional
	os string,
) *Directory {
	if arch == "" {
		arch = runtime.GOARCH
	}
	if os == "" {
		os = runtime.GOOS
	}

	if source != nil {
		g = g.WithProject(source)
	}

	command := append([]string{"go", "build"}, args...)
	return g.prepare().
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("GOOS", os).
		WithExec(command).
		Directory(PROJ_MOUNT)
}

// Build a Go project returning a Container containing the build
func (g *Golang) BuildContainer(
	// The Go source code to build
	// +optional
	source *Directory,
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
	base *Container,
) *Container {
	args = append(args, "-o", "/src/build-output/")
	dir := g.Build(source, args, arch, os)
	if base == nil {
		base = dag.Container().From("ubuntu:latest")
	}
	return base.
		WithDirectory("/usr/local/bin/", dir.Directory("./build-output/"))
}

// Test the Go project
func (g *Golang) Test(
	ctx context.Context,
	// The Go source code to test
	// +optional
	source *Directory,
	// Arguments to `go test`
	// +optional
	// +default "./..."
	args []string,
) (string, error) {
	if source != nil {
		g = g.WithProject(source)
	}
	command := append([]string{"go", "test"}, args...)
	return g.prepare().WithExec(command).Stdout(ctx)
}

// Lint the Go project
func (g *Golang) GolangciLint(
	ctx context.Context,
	// The Go source code to lint
	// +optional
	source *Directory,
) (string, error) {
	if source != nil {
		g = g.WithProject(source)
	}
	return dag.Container().From(LINT_IMAGE).
		WithMountedDirectory("/src", g.Proj).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}).
		Stdout(ctx)
}

// Sets up the Container with a golang image and cache volumes
func (g *Golang) Base(version string) *Golang {
	mod := dag.CacheVolume("gomodcache")
	build := dag.CacheVolume("gobuildcache")
	image := fmt.Sprintf("golang:%s", version)
	c := dag.Container().
		From(image).
		WithMountedCache("/go/pkg/mod", mod).
		WithMountedCache("/root/.cache/go-build", build)
	g.Ctr = c
	return g
}

// The go build container
func (g *Golang) Container() *Container {
	return g.Ctr
}

// The go project directory
func (g *Golang) Project() *Directory {
	return g.Ctr.Directory(PROJ_MOUNT)
}

// Specify the Project to use in the module
func (g *Golang) WithProject(dir *Directory) *Golang {
	g.Proj = dir
	return g
}

// Bring your own container
func (g *Golang) WithContainer(ctr *Container) *Golang {
	g.Ctr = ctr
	return g
}

// Build a remote git repo
func (g *Golang) BuildRemote(
	remote, ref, module string,
	// +optional
	arch string,
	// +optional
	platform string,
) *Directory {
	git := dag.Git(fmt.Sprintf("https://%s", remote)).
		Branch(ref).
		Tree()
	g = g.WithProject(git)

	if arch == "" {
		arch = runtime.GOARCH
	}
	if platform == "" {
		platform = runtime.GOOS
	}
	command := append([]string{"go", "build", "-o", "build/"}, module)
	return g.prepare().
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("GOOS", platform).
		WithExec(command).
		Directory(fmt.Sprintf("%s/%s/", PROJ_MOUNT, "build"))
}

// Private func to check readiness and prepare the container for build/test/lint
func (g *Golang) prepare() *Container {
	c := g.Ctr.
		WithDirectory(PROJ_MOUNT, g.Proj).
		WithWorkdir(PROJ_MOUNT)
	return c
}
