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
	Ctr  *Container
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
	args []string,
	// +optional
	arch string,
	// +optional
	os string,
	// +optional
	project *Directory,
) *Directory {
	if arch == "" {
		arch = runtime.GOARCH
	}
	if os == "" {
		os = runtime.GOOS
	}
	if project != nil {
		g.Proj = project
	}

	command := append([]string{"go", "build"}, args...)
	return g.prepare().
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("GOOS", os).
		WithExec(command).
		Directory(PROJ_MOUNT)
}

// Test the Go project
func (g *Golang) Test(ctx context.Context, args []string) (string, error) {
	command := append([]string{"go", "test"}, args...)
	return g.prepare().WithExec(command).Stdout(ctx)
}

// Lint the Go project
func (g *Golang) GolangciLint(ctx context.Context) (string, error) {
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

// Accessor for the Container
func (g *Golang) Container() *Container {
	return g.Ctr
}

// Accessor for the Project
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
