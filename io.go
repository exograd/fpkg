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
