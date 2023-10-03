package circle

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Global orb list for command matching during unmarshalling
var Glorbs = map[string]*Orb{}

type Config struct {
	version   string               `yaml:"version"`
	jobs      map[string]*Job      `yaml:"jobs"`
	workflows map[string]*Workflow `yaml:"workflows"`
	orbs      map[string]*Orb      `yaml:"orbs"`
}

type Job struct {
	docker []*Docker `yaml:"docker"`
	steps  []*Step   `yaml:"steps"`
}

type Workflow struct {
	jobs []string `yaml:"jobs"`
}

type Docker struct {
	image string `yaml:"image"`
}

// Custom parser because orbs have to be evaluated before jobs
func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	nodes := map[string]*yaml.Node{}
	for i := 0; i < len(value.Content); i += 2 {
		k := value.Content[i]
		if k.Tag == "!!str" {
			nodes[value.Content[i].Value] = value.Content[i+1]
		}
	}

	// config.Version
	if nodes["version"] != nil {
		err := nodes["version"].Decode(&c.version)
		if err != nil {
			return err
		}
	}

	// config.Orbs
	if nodes["orbs"] != nil {
		err := nodes["orbs"].Decode(&c.orbs)
		if err != nil {
			return err
		}
	}
	Glorbs = c.orbs

	// config.Jobs
	if nodes["jobs"] != nil {
		err := nodes["jobs"].Decode(&c.jobs)
		if err != nil {
			return err
		}
	}

	// config.Workflows
	if nodes["workflows"] != nil {
		err := nodes["workflows"].Decode(&c.workflows)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadConfig(path string) (*Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// parse yaml
	var configParsed *Config
	err = yaml.Unmarshal(configBytes, &configParsed)

	return configParsed, err
}
