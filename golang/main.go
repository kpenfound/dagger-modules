package main

import (
	"errors"
	"fmt"
	"context"
)

const DEFAULT_GO = "1.21"

type Golang struct {
	Ctr *Container
	Project *Directory
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
	g.Project = d
	return g
}

func (g *Golang) Build(ctx context.Context, args []string) (*Golang, error) {
	g, err := g.checkReadiness(ctx)
	if err != nil {
		return nil, err
	}
	command := append([]string{"go", "build"}, args...)
	c, err := g.Ctr.WithExec(command).Sync(ctx)
	if err != nil {
		return nil, err
	}
	g.Ctr = c
	return g, nil
}

func (g *Golang) Test(ctx context.Context, args []string) (*Golang, error) {
	g, err := g.checkReadiness(ctx)
	if err != nil {
		return nil, err
	}
	command := append([]string{"go", "test"}, args...)
	c, err := g.Ctr.WithExec(command).Sync(ctx)
	if err != nil {
		return nil, err
	}
	g.Ctr = c
	return g, nil
}

func (g *Golang) GolangciLint(ctx context.Context) (*Golang, error) {
	_, err := g.checkReadiness(ctx) // Dont override g here
	if err != nil {
		return nil, err
	}
	_, err = dag.Container().From("golangci/golangci-lint:v1.48").
		WithMountedDirectory("/src", g.Project).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Golang) checkReadiness(ctx context.Context) (*Golang, error) {
	if g.Project == nil {
		return nil, errors.New("Golang: Project is not set. Must call WithProject before executing")
	}

	if g.Ctr == nil {
		gd, err := g.Base(ctx, DEFAULT_GO)
		if err != nil {
			return nil, err
		}
		g = gd
	}
	return g, nil
}

