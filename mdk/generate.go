package main

import (
	"bytes"
	"context"
	"dagger/mdk/internal/dagger"
	"fmt"
	"text/template"

	"github.com/gobeam/stringy"
)

type Generate struct {
	Source *dagger.Directory
}

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
		if t == "OBJECT_KIND" {
			objName, err := o.AsObject().Name(ctx)
			if err != nil {
				return "", nil, err
			}
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

func templateExample(sdk string, module string, objects []string) (*dagger.Directory, error) {
	templ := ""
	path := ""
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

	return dag.Directory().WithNewFile(path, out.String()), nil
}

func toPascal(s string) string {
	str := stringy.New(s)
	return str.PascalCase().Get()
}

func toSnake(s string) string {
	str := stringy.New(s)
	return str.SnakeCase().ToLower()
}

func toCamel(s string) string {
	str := stringy.New(s)
	return str.CamelCase().Get()
}

func goTemplate() string {
	return `
// {{ .Module |  ToPascal }} examples in Go
package main

type Examples struct{}

{{ range .Objects }}
// Example for {{ . | ToPascal }} function
func {{ $.Module | ToPascal }}_{{ . | ToPascal }} () {
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
