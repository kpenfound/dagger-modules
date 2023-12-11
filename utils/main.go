package main

import (
	"context"
	"fmt"
)

type Utils struct {}

// Get a tarball of a Directory
func (m *Utils) Tar(dir *Directory) *File {
	ctr := dag.Container().
	From("alpine:3.18").
	WithMountedDirectory("/assets", dir).
	WithExec([]string{"tar", "czf", "out.tar.gz", "/assets"}).
	File("out.tar.gz")
}

func (m *Utils) Multisync(ctx context.Context, ctrs []*Container) error {
	hck := dag.Directory()

	for i, ctr := range ctrs {
		ctrFile := fmt.Sprintf("/syncfile%v", i)
		ctr = ctr.WithNewFile("/syncfile")
		hck = hck.WithFile(ctrFile, ctr.File("/syncfile"))
	}

	_, err := hck.Entries()

	return err
}
