package main

import (
	"bytes"
	"context"
	"dagger/mdk/internal/dagger"
	"fmt"
	"text/template"

	"github.com/gobeam/stringy"
)

// Utilities for generating module things
type Generate struct {
	Source *dagger.Directory
}

// Generate Daggerverse examples for a module
func (g *Generate) Examples(ctx context.Context) (*dagger.Directory, error) {
	dir := dag.Directory()

	examples := []string{"go", "python", "typescript"}

	for _, e := range examples {
		example := dag.Directory().WithNewFile("/dagger.json", fmt.Sprintf(`
{
	"name": "example",
	"sdk": "%s",
	"engineVersion": "v0.13.3"
}
`, e))

		exampleMod := example.AsModule()

		// Generate our specific functions
		name, funcs, err := g.getObjects(ctx)
		if err != nil {
			return nil, err
		}

		templateContent, err := templateExample(e, name, funcs)
		if err != nil {
			return nil, err
		}

		example = example.WithDirectory("/", exampleMod.GeneratedContextDirectory())

		example = example.WithDirectory("/", templateContent)

		dir = dir.WithDirectory("/examples/"+e, example)
	}
	return dir, nil
}

// Get the name and functions of the objects in the g.Source module
func (g *Generate) getObjects(ctx context.Context) (string, []string, error) {
	objs := []string{}
	mod := g.Source.AsModule().Initialize()
	name, err := mod.Name(ctx)
	if err != nil {
		return "", nil, err
	}

	objects, err := mod.Objects(ctx)
	if err != nil {
		return "", nil, err
	}

	for _, o := range objects {
		t, err := o.Kind(ctx)
		if err != nil {
			return "", nil, err
		}
		// Only look at Objects
		if t == "OBJECT_KIND" {
			objName, err := o.AsObject().Name(ctx)
			if err != nil {
				return "", nil, err
			}
			// Only look at objects that match the module name
			if toCamel(objName) == toCamel(name) {
				funcs, err := o.AsObject().Functions(ctx)
				if err != nil {
					return "", nil, err
				}
				for _, f := range funcs {
					funcName, err := f.Name(ctx)
					if err != nil {
						return "", nil, err
					}
					objs = append(objs, funcName)
				}
				return name, objs, nil
			}
		}

	}

	return name, objs, nil
}

// Generate a template for the given SDK
func templateExample(sdk string, module string, objects []string) (*dagger.Directory, error) {
	templ := ""
	path := ""

	// figure out which SDK we're templating
	switch sdk {
	case "go":
		templ = goTemplate()
		path = "main.go"
	case "python":
		templ = pythonTemplate()
		path = "src/main/__init__.py"
	case "typescript":
		templ = typescriptTemplate()
		path = "src/index.ts"
	}

	// Write out template
	funcMap := template.FuncMap{
		"ToPascal": toPascal,
		"ToSnake":  toSnake,
		"ToCamel":  toCamel,
	}
	t, err := template.New("examples").Funcs(funcMap).Parse(templ)
	if err != nil {
		return nil, err
	}

	templData := struct {
		Module  string
		Objects []string
	}{Module: module, Objects: objects}
	var out bytes.Buffer
	err = t.Execute(&out, templData)
	if err != nil {
		return nil, err
	}

	// return directory with generated example file
	return dag.Directory().WithNewFile(path, out.String()), nil
}

// PascalCase
func toPascal(s string) string {
	str := stringy.New(s)
	return str.PascalCase().Get()
}

// snake_case
func toSnake(s string) string {
	str := stringy.New(s)
	return str.SnakeCase().ToLower()
}

// camelCase
func toCamel(s string) string {
	str := stringy.New(s)
	return str.CamelCase().Get()
}

func goTemplate() string {
	return `
// {{ .Module |  ToPascal }} examples in Go
package main

type Example struct{}

{{ range .Objects }}
// Example for {{ . | ToPascal }} function
func (m *Example) {{ $.Module | ToPascal }}{{ . | ToPascal }} () {
	// TODO: implement example here
}

{{ end }}
`
}

func pythonTemplate() string {
	return `
"""{{ .Module |  ToSnake }} examples in Python"""
import dagger
from dagger import dag, function, object_type

@object_type
class Example:
{{ range .Objects }}
	@function
	def {{ $.Module | ToSnake }}_{{ . | ToSnake }}(self):
		"""Example for {{ . | ToSnake }} function"""

{{ end }}
`
}

func typescriptTemplate() string {
	return `
// {{ .Module |  ToCamel }} examples in TypeScript
import { dag, object, func } from "@dagger.io/dagger";

@object()
class Example {
{{ range .Objects }}
	/**
	 * example for {{ . | ToCamel }} function
	 */
	@func()
	{{ $.Module | ToCamel }}{{ . | ToPascal }}() {
		// TODO: implement example here
	}
{{ end }}
}
`
}

// Escape hatch for debugging
func (g *Generate) Debug() *dagger.Container {
	return dag.Container().From("alpine").WithMountedDirectory("/src", g.Source).WithWorkdir("/src")
}
