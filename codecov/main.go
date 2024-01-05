package main

import (
	"context"
	"strings"
)
type Codecov struct{}

const (
	UPLOADER = "https://uploader.codecov.io/v0.7.1/linux/codecov"
)

// upload: upload coverage reports to codecov.io
func (c *Codecov) Upload(ctx context.Context, dir *Directory, token *Secret, name Optional[string], verbose Optional[bool], files Optional[[]string], flags Optional[[]string]) (string, error) {
	command := []string{"/bin/codecov", "-t", "$CODECOV_TOKEN"}

	name_, isset := name.Get()
	if isset {
		command = append(command, "-n", name_)
	}

	verbose_ := verbose.GetOr(false)
	if verbose_ {
		command = append(command, "-v")
	}

	files_, isset := files.Get()
	if isset {
		for _, f := range files_ {
			command = append(command, "-f", f)
		}
	}

	flags_, isset := flags.Get()
	if isset {
		for _, f := range flags_ {
			command = append(command, f)
		}
	}

	return dag.Container(ContainerOpts{ Platform: "linux/amd64" }).
		From("cgr.dev/chainguard/wolfi-base").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "-o", "/bin/codecov", "-s", UPLOADER}). // TODO: validate uploader
		WithExec([]string{"chmod", "+x", "/bin/codecov"}).
		WithExec([]string{"ls", "-lah", "/bin/codecov"}).
		WithMountedDirectory("/src", dir).
		WithWorkdir("/src").
		WithSecretVariable("CODECOV_TOKEN", token).
		WithExec([]string{"sh", "-c", strings.Join(command, " ")}).Stdout(ctx)
}
