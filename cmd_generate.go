package main

import (
	"fmt"
	"path"

	"github.com/exograd/go-program"
)

func cmdGenerate(p *program.Program) {
	dirPath := p.ArgumentValue("directory")
	if dirPath == "" {
		dirPath = "."
	}

	configPath := p.OptionValue("config")
	var config GenerationConfig
	if err := config.LoadFile(configPath); err != nil {
		p.Fatal("cannot load configuration file from %s: %v", configPath, err)
	}

	if err := generateManifest(&config, dirPath); err != nil {
		p.Fatal("cannot generate manifest: %v", err)
	}

	// TODO Generate +PRE_INSTALL

	// TODO Create the archive
}

func generateManifest(config *GenerationConfig, dirPath string) error {
	var m Manifest

	m.Name = config.Name
	m.Version = config.Version
	m.Comment = config.ShortDescription
	m.Desc = config.LongDescription
	m.WWW = config.WebsiteURI
	m.Maintainer = config.Maintainer
	m.Arch = config.Architecture

	m.Deps = make(ManifestDeps, len(config.Dependencies))
	for _, dep := range config.Dependencies {
		m.Deps[dep.Name] = ManifestDep{
			Origin:  dep.Name,
			Version: dep.Version,
		}
	}

	m.Users = make([]string, len(config.Users))
	for i, user := range config.Users {
		m.Users[i] = user.Name
	}

	m.Groups = make([]string, len(config.Groups))
	for i, group := range config.Groups {
		m.Groups[i] = group.Name
	}

	m.Prefix = "/"

	// TODO files

	// TODO directories

	filePath := path.Join(dirPath, "+MANIFEST")
	if err := m.WriteFile(filePath); err != nil {
		return fmt.Errorf("cannot write %s: %w", filePath, err)
	}

	return nil
}
