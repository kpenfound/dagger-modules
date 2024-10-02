// The unofficial MDK (Module Developer Kit)
//
// The MDK is the Module Developer Kit. It provides
// utilities for develping modules.

package main

import "dagger/mdk/internal/dagger"

type Mdk struct {
	Generate Generate
}

func New(
	// +defaultPath="/"
	source *dagger.Directory,
) *Mdk {
	return &Mdk{
		Generate: Generate{Source: source},
	}
}
