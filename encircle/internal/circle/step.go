package circle

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kpenfound/dagger-modules/encircle"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type Step struct {
	name    string `yaml:"name"`
	run     *Run   `yaml:"run"`
	command *OrbCommandExecution
	workDir string `yaml:"working_directory"`
}

type Run struct {
	name        string            `yaml:"name"`
	command     string            `yaml:"command"`
	environment map[string]string `yaml:"environment"`
}

type OrbCommandExecution struct {
	orbCommand *OrbCommand
	parameters map[string]string
}

func (s *Step) toDagger(c *encircle.Container, params map[string]string) *encircle.Container {
	c = c.Pipeline(s.name)
	if s.workDir != "" { // workdir relative to project root
		c = c.WithWorkdir(filepath.Join("/src", s.workDir))
	}
	if s.run != nil {
		c = s.run.toDagger(c, params)
	}
	if s.command != nil {
		// Get default params
		maps.Copy(params, s.command.orbCommand.getDefaultParameters())
		// Override user params
		maps.Copy(params, s.command.parameters)
		c = s.command.orbCommand.toDagger(c, params)
	}

	return c
}

func (r *Run) toDagger(c *encircle.Container, params map[string]string) *encircle.Container {
	c = c.Pipeline(replaceParams(r.name, params))
	// Set env vars
	for k, v := range r.environment {
		c = c.WithEnvVariable(k, replaceParams(v, params))
	}
	// Exec command
	command := replaceParams(r.command, params)
	command = fmt.Sprintf("#!/bin/bash\n%s", command)
	script := fmt.Sprintf("/%s.sh", getSha(command))
	c = c.WithNewFile(script, encircle.ContainerWithNewFileOpts{
		Permissions: 0777,
		Contents:    command,
	})
	c = c.WithExec([]string{script})
	// Unset env vars
	for k := range r.environment {
		c = c.WithoutEnvVariable(k)
	}
	return c
}

func (s *Step) UnmarshalYAML(value *yaml.Node) error {
	switch value.Tag {
	case "!!str": // Basic command like checkout
		if value.Value == "checkout" {
			fmt.Println("warning: skipping checkout for local dev")
		} else if strings.Contains(value.Content[0].Value, "/") {
			// Handle orb command with no params
			commandParts := strings.Split(value.Content[0].Value, "/")
			orb := commandParts[0]
			command := commandParts[1]
			s.command = &OrbCommandExecution{
				orbCommand: findCommandForOrb(orb, command),
			}
		} else {
			fmt.Printf("warning: unknown step command: %s\n", value.Value)
		}
	case "!!map":
		if len(value.Content) == 0 {
			break
		}
		// Basic run block
		if value.Content[0].Value == "run" {
			// run block
			if value.Content[1].Tag == "!!map" {
				var r *Run
				err := value.Content[1].Decode(&r)
				if err != nil {
					return err
				}
				s.run = r
			}
			// inline run
			if value.Content[1].Tag == "!!str" {
				r := &Run{}
				r.command = value.Content[1].Value
				s.run = r
			}

			// handle orb command with params
		} else if strings.Contains(value.Content[0].Value, "/") {
			commandParts := strings.Split(value.Content[0].Value, "/")
			orb := commandParts[0]
			command := commandParts[1]
			s.command = &OrbCommandExecution{
				orbCommand: findCommandForOrb(orb, command),
			}

			// parse params
			if len(value.Content) > 1 {
				params := map[string]string{}
				for i := 0; i < len(value.Content[1].Content); i += 2 {
					k := value.Content[1].Content[i].Value
					v := value.Content[1].Content[i+1].Value
					params[k] = v
				}
				s.command.parameters = params
			}
		} else {
			fmt.Printf("warning: unhandled command %s\n", value.Content[0].Value)
		}
	default:
		fmt.Printf("Unknown yaml Tag %s\n", value.Tag)
	}
	return nil
}

func findCommandForOrb(orb string, command string) *OrbCommand {
	if Glorbs[orb] != nil {
		return Glorbs[orb].orb.commands[command]
	}
	fmt.Println("didnt find orb")
	return nil
}

func replaceParams(target string, params map[string]string) string {
	if strings.Contains(target, "<< parameters.") || strings.Contains(target, "<<parameters.") {
		for k, v := range params {
			p1 := fmt.Sprintf("<< parameters.%s >>", k)
			p2 := fmt.Sprintf("<<parameters.%s>>", k)
			target = strings.ReplaceAll(target, p1, v)
			target = strings.ReplaceAll(target, p2, v)
		}
	}
	return target
}

func getSha(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
