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
	// +optional
	name string, // optional name
	// +optional
	verbose bool, // optional verbose output
	// +optional
	files []string, // optional list of coverage files
	// +optional
	flags []string, // optional additional flags for uploader
) (string, error) {
	command := []string{"/bin/codecov", "-t", "$CODECOV_TOKEN"}

	if name != "" {
		command = append(command, "-n", name)
	}

	if verbose {
		command = append(command, "-v")
	}

	if len(files) > 0 {
		for _, f := range files {
			command = append(command, "-f", f)
		}
	}

	if len(flags) > 0 {
		for _, f := range flags {
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
