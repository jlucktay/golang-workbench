package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func getResponse(get url.URL) (io.Reader, error) {
	req, reqErr := http.NewRequest("GET", get.String(), nil)
	if reqErr != nil {
		Error.Printf("URL '%s': request error: %v\n", get.String(), reqErr)
		return nil, reqErr
	}

	req.Header.Add("User-Agent", "jlucktay (monzo-crawler)")
	res, resErr := http.DefaultClient.Do(req)

	// HTTP response errors and non-200 status codes will
	// 1) log to an error file, and
	// 2) return a nil buffer and the error from http.DefaultClient
	if resErr != nil {
		Error.Printf("URL '%s': response error: %v\n", get.String(), resErr)
		return nil, resErr
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("URL '%s': status code error: [%d] %s",
			get.String(), res.StatusCode, res.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	return buf, nil
}
