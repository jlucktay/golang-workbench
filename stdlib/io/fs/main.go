package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const specPathBase = `/Users/jameslucktaylor/git/github.com/TykTechnologies/ara/k8s/deployments/home/go`

func main() {
	sfs := SpecFS{base: specPathBase}
	if err := StatHomeNS(sfs); err != nil {
		fmt.Fprintf(os.Stderr, "could not stat file: %v\n", err)
	}
}

func StatHomeNS(f fs.FS) error {
	file, err := f.Open("home_namespace.yaml")
	if err != nil {
		return err
	}

	s, err := file.Stat()
	if err != nil {
		return err
	}

	fmt.Printf("stat: '%+v'\n", s)

	return nil
}

type SpecFS struct {
	base string
}

func (s SpecFS) Open(name string) (fs.File, error) {
	o := filepath.Join(s.base, name)

	f, err := os.Open(o)
	if err != nil {
		return nil, err
	}

	return f, nil
}
