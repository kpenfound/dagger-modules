package circle

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Step struct {
	Name    string `yaml:"name"`
	Run     *Run   `yaml:"run"`
	Command *OrbCommandExecution
	WorkDir string `yaml:"working_directory"`
}

type Run struct {
	Name        string            `yaml:"name"`
	Command     string            `yaml:"command"`
	Environment map[string]string `yaml:"environment"`
}

type OrbCommandExecution struct {
	OrbCommand *OrbCommand
	Parameters map[string]string
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
			s.Command = &OrbCommandExecution{
				OrbCommand: findCommandForOrb(orb, command),
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
				s.Run = r
			}
			// inline run
			if value.Content[1].Tag == "!!str" {
				r := &Run{}
				r.Command = value.Content[1].Value
				s.Run = r
			}

			// handle orb command with params
		} else if strings.Contains(value.Content[0].Value, "/") {
			commandParts := strings.Split(value.Content[0].Value, "/")
			orb := commandParts[0]
			command := commandParts[1]
			s.Command = &OrbCommandExecution{
				OrbCommand: findCommandForOrb(orb, command),
			}

			// parse params
			if len(value.Content) > 1 {
				params := map[string]string{}
				for i := 0; i < len(value.Content[1].Content); i += 2 {
					k := value.Content[1].Content[i].Value
					v := value.Content[1].Content[i+1].Value
					params[k] = v
				}
				s.Command.Parameters = params
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
		return Glorbs[orb].Orb.Commands[command]
	}
	fmt.Println("didnt find orb")
	return nil
}
