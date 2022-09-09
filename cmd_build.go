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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
	"unicode"

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

	manifest, err := generateManifest(config, dirPath)
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

	if err := createArchive(config, dirPath, manifest, archive); err != nil {
		if removeErr := os.Remove(archivePath); removeErr != nil {
			p.Error("cannot delete %q: %v", archivePath, removeErr)
		}

		p.Fatal("cannot create archive: %v", err)
	}

	fmt.Printf("%s\n", archivePath)
}

func generateManifest(config *GenerationConfig, dirPath string) (*Manifest, error) {
	m := NewManifest()

	m.Name = config.Name
	m.Version = config.Version
	m.Comment = config.ShortDescription
	m.Desc = config.LongDescription
	m.WWW = config.WebsiteURI
	m.Maintainer = config.Maintainer

	if arch := config.Architecture; arch != "" {
		m.Arch = config.Architecture
	} else {
		// "pkg add" will segfault if there is no arch field. See
		// https://github.com/freebsd/pkg/issues/2070.
		m.Arch = "*"
	}

	if longDesc := config.LongDescription; longDesc != "" {
		m.Desc = longDesc
	} else {
		desc := []rune(m.Comment)
		desc[0] = unicode.ToUpper(desc[0])
		m.Desc = string(desc) + "."
	}

	if origin := config.Origin; origin != "" {
		m.Origin = origin
	} else {
		m.Origin = "misc/" + config.Name
	}

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

	err := WalkDir(dirPath, func(relPath string, info fs.FileInfo) error {
		fullPath := path.Join(dirPath, relPath)

		if !info.Mode().IsRegular() {
			return nil
		}

		checksum, err := FileSHA256Checksum(fullPath)
		if err != nil {
			return fmt.Errorf("cannot compute checksum of %q: %w",
				fullPath, err)
		}

		permString := strconv.FormatInt(int64(info.Mode().Perm()), 8)

		if info.IsDir() {
			m.Directories[relPath] = ManifestDirectory{
				Uname: config.FileOwner,
				Gname: config.FileGroup,
				Perm:  permString,
			}
		} else {
			m.Files[relPath] = ManifestFile{
				Uname: config.FileOwner,
				Gname: config.FileGroup,
				Perm:  permString,
				Sum:   hex.EncodeToString(checksum),
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func createArchive(config *GenerationConfig, dirPath string, manifest *Manifest, archive io.Writer) error {
	now := time.Now().UTC()

	w := tar.NewWriter(archive)

	addFile := func(name string, mode int64, owner, group string, data []byte) error {
		header := tar.Header{
			Typeflag: tar.TypeReg,
			Name:     name,
			Size:     int64(len(data)),
			Mode:     mode,
			ModTime:  now,
			Uname:    owner,
			Gname:    group,
		}

		if data == nil {
			header.Typeflag = tar.TypeDir
		}

		if owner := config.FileOwner; owner != "" {
			header.Uname = owner
		}

		if group := config.FileGroup; group != "" {
			header.Gname = group
		}

		if err := w.WriteHeader(&header); err != nil {
			return fmt.Errorf("cannot write header: %w", err)
		}

		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("cannot write data: %w", err)
		}

		return nil
	}

	manifestData, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("cannot encode manifest: %w", err)
	}

	err = addFile("+MANIFEST", 0644, config.FileOwner, config.FileGroup,
		manifestData)
	if err != nil {
		return fmt.Errorf("cannot add manifest: %w", err)
	}

	// TODO +PRE_INSTALL

	relPaths := make([]string, 0, len(manifest.Files))
	for relPath := range manifest.Files {
		relPaths = append(relPaths, relPath)
	}
	sort.Strings(relPaths)

	for _, relPath := range relPaths {
		mfile := manifest.Files[relPath]
		filePath := path.Join(dirPath, relPath)

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("cannot read %q: %w", filePath, err)
		}

		perm, err := strconv.ParseInt(mfile.Perm, 8, 64)
		if err != nil {
			return fmt.Errorf("cannot parse permission string %q: %w",
				mfile.Perm, err)
		}

		err = addFile(relPath, perm, mfile.Uname, mfile.Gname, data)
		if err != nil {
			return fmt.Errorf("cannot add %q: %w", filePath, err)
		}
	}

	relPaths = make([]string, 0, len(manifest.Directories))
	for relPath := range manifest.Directories {
		relPaths = append(relPaths, relPath)
	}
	sort.Strings(relPaths)

	for _, relPath := range relPaths {
		mdir := manifest.Directories[relPath]
		filePath := path.Join(dirPath, relPath)

		perm, err := strconv.ParseInt(mdir.Perm, 8, 64)
		if err != nil {
			return fmt.Errorf("cannot parse permission string %q: %w",
				mdir.Perm, err)
		}

		err = addFile(relPath, perm, mdir.Uname, mdir.Gname, nil)
		if err != nil {
			return fmt.Errorf("cannot add directory %q: %w", filePath, err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("cannot close archive: %w", err)
	}

	return nil
}
