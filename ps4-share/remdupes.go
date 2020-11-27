package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeebo/xxh3"
)

// - go to each directory
//   - get all files in the directory
//     - are any the same size?
//       - any identical hashes?

const pathToWalk = "/Volumes/Sgte-ExFAT/PS4/SHARE/"

func main() {
	files := map[string][]os.FileInfo{}

	err := filepath.Walk(pathToWalk, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failure accessing path %q: %w\n", fullPath, err)
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		dir := strings.TrimSuffix(fullPath, info.Name())

		if files[dir] == nil {
			files[dir] = []os.FileInfo{}
		}

		files[dir] = append(files[dir], info)

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking path: %v\n", err)

		return
	}

	for dir, ls := range files {
	OuterFileLoop:
		for outer, fi1 := range ls {
		InnerFileLoop:
			for _, fi2 := range ls[outer+1:] {
				if fi1.Size() != fi2.Size() {
					continue
				}

				file1, errRead1 := ioutil.ReadFile(dir + fi1.Name())
				if errRead1 != nil {
					fmt.Fprintf(os.Stderr, "could not read '%s': %v", dir+fi1.Name(), errRead1)

					continue OuterFileLoop
				}

				hash1 := xxh3.Hash(file1)

				file2, errRead2 := ioutil.ReadFile(dir + fi2.Name())
				if errRead2 != nil {
					fmt.Fprintf(os.Stderr, "could not read '%s': %v", dir+fi2.Name(), errRead2)

					continue InnerFileLoop
				}

				hash2 := xxh3.Hash(file2)

				if hash1 == hash2 {
					fmt.Printf("hash match!\n\n")
				}
			}
		}
	}
}
