package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
)

const (
	// exitFail is the exit code if the program fails.
	exitFail = 1
)

const (
	baseURL = "https://www.giantbomb.com/api/videos/"
)

// URL: https://www.giantbomb.com/api/videos/?api_key=[YOUR API KEY]

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	if errConf := gatherConfig(); errConf != nil {
		return errConf
	}

	fmt.Printf("api key: '%s'\n", viper.GetString("api-key"))

	return nil
}
