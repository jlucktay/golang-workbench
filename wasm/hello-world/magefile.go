//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// Default denotes Mage's default target when invoked without one explicitly.
	Default = Def

	// Aliases links short versions to longer target names.
	Aliases = map[string]interface{}{
		"b": Build,
		"c": Clean,
		"s": Serve,
	}
)

// Def cleans, builds, serves, in that order.
func Def() {
	mg.Deps(Clean)
	mg.Deps(Build)
	mg.Deps(Serve)
}

// Build the web app into webassembly.
func Build() error {
	return sh.RunWith(
		map[string]string{
			"GOARCH": "wasm",
			"GOOS":   "js",
		},
		"go",
		"build",
		"-a",
		"-o",
		"content/lib.wasm",
		"content/content.go",
	)
}

// Clean will delete various bits of cruft.
func Clean() error {
	return sh.RunV("rm", "-fv", "content/lib.wasm")
}

// Serve will run up the server locally.
func Serve() error {
	return sh.RunV("go", "run", "server/server.go", "--dir=content")
}
