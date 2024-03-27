// A module for working with Postgres
//
// Module for building Postgres containers, serving a postgres service,
// or querying a running postgres database.

package main

import (
	"fmt"
)

type Postgres struct{
	Version string
}

func New(
	// The postgres container tag
	// +optional
	// +default="16"
	version string,
) *Postgres {
	return &Postgres{
		Version: version,
	}
}

// Specify a tag of postgres to use
func (n *Postgres) WithVersion(
	// The postgres container tag
	version string,
) *Postgres {
	n.Version = version
	return n
}

// Get the postgres container with the directory in it
func (n *Postgres) Container(
	// The platform of the container
	// +optional
	platform string,
) *Container {
	opts := ContainerOpts{}
	if platform != "" {
		opts.Platform = Platform(platform)
	}

	return dag.Container(opts).
		From("postgres:" + n.Version).
		WithExposedPort(5432)
}

// Get a postgres container as a service
func (n *Postgres) Service(
	// The platform of the container
	// +optional
	platform string,
) *Service {
	return n.Container(platform).AsService()
}

// Get a terminal of a psql client connected to a Postgres database
func (n *Postgres) Client(
	// The Postgres server to connect to
	server *Service,
	// The postgres database
	db string,
	// The postgres user
	user string,
	// The postgres password
	password string,
	// The postgres port
	// +optional
	// +default="5432"
	port string,
	// The postgres client container tag
	// +optional
	// +default="16"
	version string,
) *Terminal {
	connectionString := fmt.Sprintf("postgresql://%s:%s@postgres:%s/%s", user, password, port, db)
	return dag.Container().
		From("postgres:" + version).
		WithoutEntrypoint().
		WithServiceBinding("postgres", server).
		WithDefaultTerminalCmd([]string{"psql", connectionString}).
		Terminal()
}

