// Copyright (c) 2022 Exograd SAS.
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
// SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR
// IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

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
	FileOwner        string                       `yaml:"file_owner,omitempty"`
	FileGroup        string                       `yaml:"file_group,omitempty"`
}

type GenerationConfigDependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type GenerationConfigUser struct {
	Name string `yaml:"name"`
	UID  uint   `yaml:"uid"`
}

type GenerationConfigGroup struct {
	Name string `yaml:"name"`
	GID  uint   `yaml:"gid"`
}

func DefaultGenerationConfig() *GenerationConfig {
	return &GenerationConfig{
		FileOwner: "root",
		FileGroup: "wheel",
	}
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
