//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	tag = "hello-docker-main"
)

var (
	// Default denotes Mage's default target when invoked without one explicitly
	Default = All

	// Aliases links short versions to longer target names
	Aliases = map[string]interface{}{
		"a": All,
		"b": Build,
		"s": Show,
		"r": Run,
		"l": Lint,
		"c": Clean,
	}

	versions = [3]string{"1.0", "1.1", "1.2"}
)

// All lints, builds, and runs, in that order
func All() {
	mg.Deps(Lint)
	mg.Deps(Build)
	mg.Deps(Run)
}

// Build the Docker images
func Build() {
	for _, version := range versions {
		sh.RunV("docker", "build", "--tag", tag+":"+version, "-f", version+"/Dockerfile", ".")
	}

	mg.Deps(Show)
}

// Show will list the images built with our tag
func Show() error {
	return sh.RunV("docker", "images", tag)
}

// Run will execute all versions of our tagged image(s)
func Run() {
	for _, version := range versions {
		sh.RunV("docker", "run", tag+":"+version)
	}
}

// Lint will check the Dockerfiles for errors
func Lint() {
	for _, version := range versions {
		sh.RunV("hadolint", version+"/Dockerfile")
	}
}

// Clean will delete all versions of the Docker image(s)
func Clean() {
	for _, version := range versions {
		sh.RunV("docker", "image", "rm", "--force", tag+":"+version)
	}
}
