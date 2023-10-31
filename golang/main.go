package main

import (
	"fmt"
	"context"
	"runtime"
)

const (
	DEFAULT_GO = "1.21"
	PROJ_MOUNT = "/src"
)

type Golang struct {
	Ctr *Container
	Proj *Directory
}

// Accessor for the Container
func (g *Golang) Container() *Container {
	return g.Ctr
}

// Accessor for the Project
func (g *Golang) Project() *Directory {
	return g.Ctr.Directory(PROJ_MOUNT)
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

// Specify the Project to use in the module
func (g *Golang) WithProject(d *Directory) (*Golang) {
	g.Proj = d
	return g
}

// Bring your own container
func (g *Golang) WithContainer(c *Container) (*Golang) {
	g.Ctr = c
	return g
}

// Build the project
func (g *Golang) Build(args []string, arch Optional[string]) *Directory {
	archStr := arch.GetOr(runtime.GOARCH)
	command := append([]string{"go", "build"}, args...)
	return g.prepare().
	WithEnvVariable("GOARCH", archStr).
	WithExec(command).
	Directory(PROJ_MOUNT)
}

// Test the project
func (g *Golang) Test(ctx context.Context, args []string) (string, error) {
	command := append([]string{"go", "test"}, args...)
	return g.prepare().WithExec(command).Stdout(ctx)
}

// Lint the project
func (g *Golang) GolangciLint(ctx context.Context) (string, error) {
	return dag.Container().From("golangci/golangci-lint:v1.48").
		WithMountedDirectory("/src", g.Proj).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}).
		Stdout(ctx)
}

// Private func to check readiness and prepare the container for build/test/lint
func (g *Golang) prepare() *Container {
	if g.Proj == nil {
		g.Proj = dag.Directory() // Unsure about this. Maybe want to error
	}

	if g.Ctr == nil {
		gd := g.Base(DEFAULT_GO)
		g.Ctr = gd.Ctr
	}

	c := g.Ctr.
	WithDirectory(PROJ_MOUNT, g.Proj).
	WithWorkdir(PROJ_MOUNT)
	return c
}

