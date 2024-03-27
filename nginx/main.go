// A module for building Nginx containers
//
// Module for building Nginx containers to serve the given
// directory.

package main

type Nginx struct{
	Version string
	Dir *Directory
}

func New(
	// The nginx container tag
	// +optional
	// +default="1.25"
	version string,
) *Nginx {
	return &Nginx{
		Version: version,
	}
}

// Specify a tag of nginx to use
func (n *Nginx) WithVersion(
	// The nginx container tag
	version string,
) *Nginx {
	n.Version = version
	return n
}

// Add a directory for nginx to serve
func (n *Nginx) WithDirectory(
	// The directory for nginx to serve
	directory *Directory,
) *Nginx {
	n.Dir = directory
	return n
}

// Get the nginx container with the directory in it
func (n *Nginx) Container(
	// The platform of the container
	// +optional
	platform string,
) *Container {
	opts := ContainerOpts{}
	if platform != "" {
		opts.Platform = Platform(platform)
	}
	ctr := dag.Container(opts).
		From("nginx:" + n.Version).
		WithExposedPort(80)

	if n.Dir != nil {
		ctr = ctr.WithDirectory("/usr/share/nginx/html", n.Dir)
	}

	return ctr
}

