package main

import (
	"fmt"
	"context"
)

type Golang struct {}

func (m *Golang) Base(ctx context.Context, version string) (*Container, error) {
	mod := dag.CacheVolume("gomodcache")
	build := dag.CacheVolume("gobuildcache")
	image := fmt.Sprintf("golang:%s", version)
	return dag.Container().
	From(image).
	WithMountedCache("/go/pkg/mod", mod).
	WithMountedCache("/root/.cache/go-build", build).
	Sync(ctx)
}

func (g *Golang) Build(ctx context.Context, c *Container, args []string) (*Container, error) {
	command := append([]string{"go", "build"}, args...)
	return c.WithExec(command).Sync(ctx)
}

func (g *Golang) Test(ctx context.Context, c *Container, args []string) (*Container, error) {
	command := append([]string{"go", "test"}, args...)
	return c.WithExec(command).Sync(ctx)
}

func (g *Golang) GolangciLint(ctx context.Context, d *Directory) (string, error) {
	return dag.Container().From("golangci/golangci-lint:v1.48").
		WithMountedDirectory("/src", d).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}).
		Stdout(ctx)
}

