// A generated module for Mdk functions

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

func (m *Mdk) Foo() string {
	return "foo"
}

func (m *Mdk) BarBaz() string {
	return "bar"
}
