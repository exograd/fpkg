package main

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type GenerationConfig struct {
	PackageName            string   `yaml:"package_name"`
	Version                string   `yaml:"version,omitempty"`
	ShortDescription       string   `yaml:"short_description,omitempty"`
	LongDescription        string   `yaml:"long_description,omitempty"`
	WebsiteURI             string   `yaml:"website_uri,omitempty"`
	MaintainerEmailAddress string   `yaml:"maintainer_email_address,omitempty"`
	Architecture           string   `yaml:"architecture,omitempty"`
	Dependencies           []string `yaml:"dependencies,omitempty"`
	Users                  []string `yaml:"users,omitempty"`
	Groups                 []string `yaml:"groups,omitempty"`
}

func (c *GenerationConfig) UnmarshalYAML(value *yaml.Node) error {
	type GenerationConfig2 GenerationConfig
	c2 := GenerationConfig2(*c)

	if err := value.Decode(&c2); err != nil {
		return err
	}

	if c2.PackageName == "" {
		return fmt.Errorf("missing or empty package name")
	}

	*c = GenerationConfig(c2)
	return nil
}

func (c *GenerationConfig) LoadFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read %q: %w", filePath, err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(&c); err != nil {
		return fmt.Errorf("cannot parse %q: %w", filePath, err)
	}

	return nil
}
