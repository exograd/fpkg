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
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/exograd/go-program"
)

func cmdBuild(p *program.Program) {
	dirPath := p.ArgumentValue("directory")
	if dirPath == "" {
		dirPath = "."
	}

	configPath := p.OptionValue("config")
	config := DefaultGenerationConfig()
	if err := config.LoadFile(configPath); err != nil {
		p.Fatal("cannot load configuration file from %s: %v", configPath, err)
	}

	if p.IsOptionSet("version") {
		config.Version = p.OptionValue("version")
	}

	if config.Version == "" {
		p.Fatal("missing or empty version")
	}

	manifest, err := generateManifest(config)
	if err != nil {
		p.Fatal("cannot generate manifest: %v", err)
	}

	// TODO Generate +PRE_INSTALL

	archivePath := manifest.PackageFilename()

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	archive, err := os.OpenFile(archivePath, flags, 0644)
	if err != nil {
		p.Fatal("cannot open %q: %v", archive, err)
	}

	if err := createArchive(config, manifest, archive); err != nil {
		if removeErr := os.Remove(archivePath); removeErr != nil {
			p.Error("cannot delete %q: %v", archivePath, removeErr)
		}

		p.Fatal("cannot create archive: %v", err)
	}

	fmt.Printf("%s\n", archivePath)
}

func generateManifest(config *GenerationConfig) (*Manifest, error) {
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

	return &m, nil
}

func createArchive(config *GenerationConfig, manifest *Manifest, archive io.Writer) error {
	now := time.Now().UTC()

	w := tar.NewWriter(archive)

	var createErr error

	addFile := func(name string, data []byte) {
		if createErr != nil {
			return
		}

		header := tar.Header{
			Typeflag: tar.TypeReg,
			Name:     name,
			Size:     int64(len(data)),
			Mode:     int64(0644),
			ModTime:  now,
		}

		if owner := config.FileOwner; owner != "" {
			header.Uname = owner
		}

		if group := config.FileGroup; group != "" {
			header.Gname = group
		}

		if err := w.WriteHeader(&header); err != nil {
			createErr = fmt.Errorf("cannot write header: %w", err)
			return
		}

		if _, err := w.Write(data); err != nil {
			createErr = fmt.Errorf("cannot write content: %w", err)
			return
		}
	}

	manifestData, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("cannot encode manifest: %w", err)
	}

	addFile("+MANIFEST", manifestData)
	// TODO +PRE_INSTALL

	// TODO files and directories

	if createErr != nil {
		return createErr
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("cannot close archive: %w", err)
	}

	return nil
}
