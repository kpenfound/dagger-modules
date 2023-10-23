package main

import (
	"errors"
	"fmt"
	"context"
)

const (
	DEFAULT_GO = "1.21"
	PROJ_MOUNT = "/src"
)

type Golang struct {
	Ctr *Container
	Proj *Directory
}

func (g *Golang) Container() *Container {
	return g.Ctr
}
func (g *Golang) Project() *Directory {
	return g.Ctr.Directory(PROJ_MOUNT)
}

func (g *Golang) Base(ctx context.Context, version string) (*Golang, error) {
	mod := dag.CacheVolume("gomodcache")
	build := dag.CacheVolume("gobuildcache")
	image := fmt.Sprintf("golang:%s", version)
	c, err := dag.Container().
	From(image).
	WithMountedCache("/go/pkg/mod", mod).
	WithMountedCache("/root/.cache/go-build", build).
	Sync(ctx)
	if err != nil {
		return nil, err
	}
	g.Ctr = c
	return g, nil
}

func (g *Golang) WithProject(d *Directory) (*Golang) {
	g.Proj = d
	return g
}

func (g *Golang) Build(ctx context.Context, args []string) (*Golang, error) {
	c, err := g.prepare(ctx)
	if err != nil {
		return nil, err
	}
	command := append([]string{"go", "build"}, args...)
	c, err = c.WithExec(command).Sync(ctx)
	if err != nil {
		return nil, err
	}
	g.Ctr = c
	return g, nil
}

func (g *Golang) Test(ctx context.Context, args []string) (*Golang, error) {
	c, err := g.prepare(ctx)
	if err != nil {
		return nil, err
	}
	command := append([]string{"go", "test"}, args...)
	c, err = c.WithExec(command).Sync(ctx)
	if err != nil {
		return nil, err
	}
	g.Ctr = c
	return g, nil
}

func (g *Golang) GolangciLint(ctx context.Context) (*Golang, error) {
	_, err := g.prepare(ctx)
	if err != nil {
		return nil, err
	}
	_, err = dag.Container().From("golangci/golangci-lint:v1.48").
		WithMountedDirectory("/src", g.Proj).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Golang) prepare(ctx context.Context) (*Container, error) {
	if g.Proj == nil {
		return nil, errors.New("Golang: Project is not set. Must call WithProject before executing")
	}

	if g.Ctr == nil {
		gd, err := g.Base(ctx, DEFAULT_GO)
		if err != nil {
			return nil, err
		}
		g = gd
	}

	c := g.Ctr.
	WithDirectory(PROJ_MOUNT, g.Proj).
	WithWorkdir(PROJ_MOUNT)
	return c, nil
}

