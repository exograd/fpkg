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
	"encoding/json"
	"fmt"
	"os"
)

// FreeBSD pkg manifests use the UCL format
// (https://github.com/vstakhov/libucl). Fortunately, UCL is officially fully
// compatible with JSON so we can just generate JSON data and be done with it.

// See https://github.com/freebsd/pkg/blob/master/libpkg/pkg_manifest.c

type Manifest struct {
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	Comment     string              `json:"comment,omitempty"`
	Desc        string              `json:"desc,omitempty"`
	Origin      string              `json:"origin,omitempty"`
	WWW         string              `json:"www,omitempty"`
	Maintainer  string              `json:"maintainer,omitempty"`
	Arch        string              `json:"arch,omitempty"`
	Deps        ManifestDeps        `json:"deps,omitempty"`
	Users       []string            `json:"users,omitempty"`
	Groups      []string            `json:"groups,omitempty"`
	Prefix      string              `json:"prefix,omitempty"`
	Files       ManifestFiles       `json:"files,omitempty"`
	Directories ManifestDirectories `json:"directories,omitempty"`
}

type ManifestDep struct {
	Origin  string `json:"origin"`
	Version string `json:"version"`
}

type ManifestDeps map[string]ManifestDep

type ManifestFile struct {
	Uname string `json:"uname"`
	Gname string `json:"gname"`
	Perm  string `json:"perm"`
	Sum   string `json:"sum"`
}

type ManifestFiles map[string]ManifestFile

type ManifestDirectory struct {
	Uname string `json:"uname"`
	Gname string `json:"gname"`
	Perm  string `json:"perm"`
}

type ManifestDirectories map[string]ManifestDirectory

func (m *Manifest) PackageFilename() string {
	return m.Name + "-" + m.Version + ".pkg"
}

func (m *Manifest) WriteFile(filePath string) error {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(m); err != nil {
		return fmt.Errorf("cannot encode manifest: %w", err)
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
