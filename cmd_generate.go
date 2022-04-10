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
	var config GenerationConfig
	if err := config.LoadFile(configPath); err != nil {
		p.Fatal("cannot load configuration file: %v", err)
	}

	fmt.Printf("config: %#v\n", config)

	fmt.Printf("dirPath: %q\n", dirPath)
}
