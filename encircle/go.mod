module github.com/kpenfound/dagger-modules/encircle

go 1.21.1

toolchain go1.21.2

require (
	github.com/99designs/gqlgen v0.17.31
	github.com/Khan/genqlient v0.6.0
	github.com/kpenfound/dagger-modules/encircle/circle v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.3.0
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	github.com/vektah/gqlparser/v2 v2.5.6 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kpenfound/dagger-modules/encircle/circle => ./circle
