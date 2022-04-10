package main

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type GenerationConfig struct {
	Name             string                       `yaml:"name"`
	Version          string                       `yaml:"version,omitempty"`
	ShortDescription string                       `yaml:"short_description,omitempty"`
	LongDescription  string                       `yaml:"long_description,omitempty"`
	WebsiteURI       string                       `yaml:"website_uri,omitempty"`
	Maintainer       string                       `yaml:"maintainer,omitempty"`
	Architecture     string                       `yaml:"architecture,omitempty"`
	Dependencies     []GenerationConfigDependency `yaml:"dependencies,omitempty"`
	Users            []GenerationConfigUser       `yaml:"users,omitempty"`
	Groups           []GenerationConfigGroup      `yaml:"groups,omitempty"`
}

type GenerationConfigDependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type GenerationConfigUser struct {
	Name string `yaml:"name"`
	UID  uint   `yaml:uid"`
}

type GenerationConfigGroup struct {
	Name string `yaml:"name"`
	GID  uint   `yaml:gid"`
}

func (c *GenerationConfig) UnmarshalYAML(value *yaml.Node) error {
	type GenerationConfig2 GenerationConfig
	c2 := GenerationConfig2(*c)

	if err := value.Decode(&c2); err != nil {
		return err
	}

	if c2.Name == "" {
		return fmt.Errorf("missing or empty package name")
	}

	*c = GenerationConfig(c2)
	return nil
}

func (c *GenerationConfig) LoadFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(&c); err != nil {
		return fmt.Errorf("cannot decode configuration: %w", err)
	}

	return nil
}
