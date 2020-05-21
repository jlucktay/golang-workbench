package main

import (
	"fmt"
	"os"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

// Photo by Marc Mintel on Unsplash
const backgroundImageFilename = "/Users/jameslucktaylor/Downloads/marc-mintel-1iYTusNPlSk-unsplash.jpg"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	dc := gg.NewContext(1200, 628)

	backgroundImage, err := gg.LoadImage(backgroundImageFilename)
	if err != nil {
		return errors.Wrap(err, "load background image")
	}

	backgroundImage = imaging.Fill(backgroundImage, dc.Width(), dc.Height(), imaging.Center, imaging.Lanczos)

	dc.DrawImage(backgroundImage, 0, 0)

	return nil
}
