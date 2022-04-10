package main

import (
	"github.com/exograd/go-program"
)

func main() {
	var c *program.Command

	p := program.NewProgram("freebsd-pkg",
		"utils to manipulate freebsd packages")

	c = p.AddCommand("generate", "generate a package", cmdGenerate)
	c.AddOptionalArgument("directory",
		"the directory containing files to package")
	c.AddOption("c", "config", "path", "freebsd-pkg.yaml",
		"the path of the configuration file")

	p.ParseCommandLine()
	p.Run()
}
