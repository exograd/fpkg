package main

import (
	"fmt"

	"github.com/exograd/go-program"
)

func cmdGenerate(p *program.Program) {
	dirPath := p.ArgumentValue("directory")
	if dirPath == "" {
		dirPath = "."
	}

	configPath := p.OptionValue("config")

	fmt.Printf("dirPath: %q\n", dirPath)
	fmt.Printf("configPath: %q\n", configPath)
}
