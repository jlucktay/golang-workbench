// +build mage

package mage

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
	}
)

// Def is assigned as the 'Default' target, so it builds.
func Def() {
	mg.Deps(Build)
}

// Build the web app using Hugo.
func Build() error {
	return sh.RunWith(
		map[string]string{
			"GOARCH": "wasm",
			"GOOS":   "js",
		},
		"go",
		"build",
		"-o",
		"content/lib.wasm",
		"content/content.go",
	)
}

// Clean will delete various bits of cruft.
func Clean() error {
	return sh.RunV("rm", "-fv", "content/lib.wasm")
}
