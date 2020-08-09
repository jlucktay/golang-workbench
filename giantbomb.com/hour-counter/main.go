package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	// exitFail is the exit code if the program fails.
	exitFail = iota
)

const (
	baseURL  = "https://www.giantbomb.com/api/videos/?format=json&field_list=%s"
	pageSize = 100
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

	vr, errGVR := getVideoResults("id", 0)
	if errGVR != nil {
		return errGVR
	}

	results := paraGet(stdout, (vr.NumberOfTotalResults/pageSize)+1)

	totalLength := 0

	results.Range(func(_, value interface{}) bool {
		xResults, ok := value.([]Results)
		if !ok {
			return false
		}

		for _, res := range xResults {
			totalLength += res.LengthSeconds
		}

		return true
	})

	dur, errPD := time.ParseDuration(fmt.Sprintf("%ds", totalLength))
	if errPD != nil {
		return errPD
	}

	fmt.Fprintf(stdout, "\nTotal length: %s (raw: %d seconds)\n", dur, totalLength)

	return nil
}

func getVideoResults(fieldList string, page int) (*VideosResult, error) {
	withFields := fmt.Sprintf(baseURL, fieldList)
	getURL := fmt.Sprintf("%s&api_key=%s&offset=%d", withFields, viper.GetString("api-key"), page*pageSize)

	u, errParse := url.Parse(getURL)
	if errParse != nil {
		return nil, errParse
	}

	req, errReq := http.NewRequestWithContext(context.TODO(), http.MethodGet, u.String(), nil)
	if errReq != nil {
		return nil, errReq
	}

	resp, errGet := http.DefaultClient.Do(req)
	if errGet != nil {
		return nil, errGet
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d %s", ErrResponseStatus, resp.StatusCode, resp.Status)
	}

	ret := &VideosResult{}

	if errDecode := json.NewDecoder(resp.Body).Decode(&ret); errDecode != nil {
		return nil, errDecode
	}

	return ret, nil
}

func paraGet(stdout io.Writer, pageLimit int) *sync.Map {
	wg := sync.WaitGroup{}
	results := sync.Map{}

	for page := 0; page <= pageLimit; page++ {
		wg.Add(1)

		go func(p int) {
			defer wg.Done()

			fmt.Fprintf(stdout, "/%d", p)

			got, errGVR := getVideoResults("id,name,length_seconds", p)
			if errGVR != nil {
				return
			}

			results.Store(p, got.Results)

			fmt.Fprintf(stdout, `\%d`, p)
		}(page)

		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	return &results
}
