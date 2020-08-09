package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

const (
	// exitFail is the exit code if the program fails.
	exitFail = iota
)

const (
	baseURL = "https://www.giantbomb.com/api/videos/?format=json&field_list=id,name,length_seconds"
)

var ErrResponseStatus = errors.New("status code not OK")

func main() {
	if err := Run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func Run(args []string, stdout io.Writer) error {
	if errConf := gatherConfig(args); errConf != nil {
		return errConf
	}

	getURL := fmt.Sprintf("%s&api_key=%s", baseURL, viper.GetString("api-key"))

	u, errParse := url.Parse(getURL)
	if errParse != nil {
		return errParse
	}

	req, errReq := http.NewRequestWithContext(context.TODO(), http.MethodGet, u.String(), nil)
	if errReq != nil {
		return errReq
	}

	resp, errGet := http.DefaultClient.Do(req)
	if errGet != nil {
		return errGet
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d %s", ErrResponseStatus, resp.StatusCode, resp.Status)
	}

	bod, errReadAll := ioutil.ReadAll(resp.Body)
	if errReadAll != nil {
		return errReadAll
	}

	if errWrite := ioutil.WriteFile("output.json", bod, 0600); errWrite != nil {
		return errWrite
	}

	return nil
}
