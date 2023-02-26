package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func main() {
	items := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	prompt := promptui.Select{
		Items:  items,
		Label:  "Select a day",
		Size:   len(items),
		Stdout: &bellSkipper{},
	}

	i, choice, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: '%v'\n", err)
		return
	}

	fmt.Printf("You chose #%d > %q.\n", i, choice)
}

// bellSkipper implements io.WriteCloser and skips the terminal bell character (ASCII code 7), writing the rest to
// os.Stderr. It is used to replace readline.Stdout, which is the package used by promptui to display the prompts.
//
// This is a workaround for the bell issue documented in:
// https://github.com/manifoldco/promptui/issues/49.
type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal bell character.
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}
