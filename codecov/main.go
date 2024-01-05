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
func (c *Codecov) Upload(
	ctx context.Context,
	dir *Directory, // Directory containing git repo and coverage output
	token *Secret, // codecov token for the git repo
	name Optional[string], // optional name
	verbose Optional[bool], // optional verbose output
	files Optional[[]string], // optional list of coverage files
	flags Optional[[]string], // optional additional flags for uploader
) (string, error) {
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
		WithExec([]string{"apk", "add", "curl", "git"}).
		WithExec([]string{"curl", "-o", "/bin/codecov", "-s", UPLOADER}). // TODO: validate uploader
		WithExec([]string{"chmod", "+x", "/bin/codecov"}).
		WithExec([]string{"ls", "-lah", "/bin/codecov"}).
		WithMountedDirectory("/src", dir).
		WithWorkdir("/src").
		WithSecretVariable("CODECOV_TOKEN", token).
		WithExec([]string{"sh", "-c", strings.Join(command, " ")}).Stdout(ctx)
}
