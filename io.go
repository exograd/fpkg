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
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
)

func WalkDir(dirPath string, fn func(string, fs.FileInfo) error) error {
	var walk func(string, string) error
	walk = func(currentPath, relPath string) error {
		info, err := os.Stat(currentPath)
		if err != nil {
			return fmt.Errorf("cannot stat %q: %w", currentPath, err)
		}

		isEmptyDir := false

		if info.IsDir() {
			entries, err := ioutil.ReadDir(currentPath)
			if err != nil {
				return fmt.Errorf("cannot list directory %q: %w",
					currentPath, err)
			}

			isEmptyDir = len(entries) == 0

			for _, entry := range entries {
				currentPath2 := path.Join(currentPath, entry.Name())
				relPath2 := path.Join(relPath, entry.Name())

				if err := walk(currentPath2, relPath2); err != nil {
					return err
				}
			}
		}

		if !info.IsDir() || isEmptyDir {
			if err := fn(relPath, info); err != nil {
				return err
			}
		}

		return nil
	}

	if err := walk(dirPath, "/"); err != nil {
		return err
	}

	return nil
}

func FileSHA256Checksum(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	return hash.Sum(nil), nil
}
