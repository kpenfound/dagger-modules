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
