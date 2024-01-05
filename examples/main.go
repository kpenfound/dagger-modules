package main

import (
	"context"
)

type Examples struct {}


func (m *Examples) Multisync(ctx context.Context) error {
	ctr1 := dag.Container().From("alpine").WithExec([]string{"apk", "add", "git"})
	ctr2 := dag.Container().From("alpine").WithExec([]string{"apk", "add", "curl"})
	ctr3 := dag.Container().From("alpine").WithExec([]string{"apk", "add", "sqlite"})

	_, err := dag.Utils().Multisync(ctx, []*Container{ctr1, ctr2, ctr3})
	return err
}

func (m *Examples) Codecov(ctx context.Context, token *Secret) (string, error) {
	coverage := `
mode: atomic
github.com/kpenfound/greetings-api/main.go:12.13,14.67 2 0
github.com/kpenfound/greetings-api/main.go:14.67,18.17 4 0
github.com/kpenfound/greetings-api/main.go:18.17,19.14 1 0
github.com/kpenfound/greetings-api/main.go:23.2,35.42 4 0
github.com/kpenfound/greetings-api/main.go:35.42,37.3 1 0
github.com/kpenfound/greetings-api/main.go:37.8,37.23 1 0
github.com/kpenfound/greetings-api/main.go:37.23,40.3 2 0
github.com/kpenfound/greetings-api/main.go:43.24,46.2 2 1
`
	repo := "https://github.com/kpenfound/greetings-api"

	src := dag.Git(repo, GitOpts{ KeepGitDir: true }).Branch("main").Tree()

	src = src.WithNewFile("coverage.out", coverage)

	return dag.Codecov().Upload(ctx, src, token)
}
