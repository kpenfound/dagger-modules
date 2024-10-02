// A generated module for Mdk functions

package main

type Mdk struct {
	Generate Generate
}

func New() *Mdk {
	return &Mdk{
		Generate: Generate{},
	}
}

func (m *Mdk) Foo() string {
	return "foo"
}

func (m *Mdk) BarBaz() string {
	return "bar"
}
