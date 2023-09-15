package main

import (
	"fmt"
	"context"
)

type Golang struct {}

func (m *Golang) Base(ctx context.Context, version string) (*Container, error) {
	cache := dag.CacheVolume("gomodcache")
	image := fmt.Sprintf("golang:%s", version)
	return dag.Container().
	From(image).
	WithMountedCache("/gomodcache", cache).
	WithEnvVariable("GOMODCACHE", "/gomodcache").Sync(ctx)
}

func (c *Container) GoBuild(ctx context.Context, args []string) (*Container, error) {
	command := append([]string{"go", "build"}, args...)
	return c.WithExec(command).Sync(ctx)
}

func (c *Container) GoTest(ctx context.Context, args []string) (*Container, error) {
	command := append([]string{"go", "test"}, args...)
	return c.WithExec(command).Sync(ctx)
}

func (d *Directory) GoTestToo(ctx context.Context, args []string) (*Container, error) {
	command := append([]string{"go", "test"}, args...)
	c, err := (&Golang{}).Base(ctx, "latest")
	if err != nil {
		return nil, err
	}
	return c.WithExec(command).Sync(ctx)
}

func (d *Directory) GoLint(ctx context.Context) (string, error) {
	return dag.Container().From("golangci/golangci-lint:v1.48").
		WithMountedDirectory("/src", d).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}).
		Stdout(ctx)
}

