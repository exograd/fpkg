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
	"github.com/exograd/go-program"
)

func main() {
	var c *program.Command

	p := program.NewProgram("fpkg",
		"tools to manipulate freebsd packages")

	c = p.AddCommand("build", "build a package", cmdBuild)
	c.AddOptionalArgument("directory",
		"the directory containing files to package")
	c.AddOption("c", "config", "path", "fpkg.yaml",
		"the path of the configuration file")
	c.AddOption("v", "version", "string", "",
		"set the version of the package")

	p.ParseCommandLine()
	p.Run()
}
