package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func getResponse(get url.URL) (io.Reader, error) {
	req, errReq := http.NewRequest("GET", get.String(), nil)
	if errReq != nil {
		Error.Printf("URL '%s': request error: %v\n", get.String(), errReq)
		return nil, errReq
	}

	req.Header.Add("User-Agent", "jlucktay (monzo-crawler)")
	res, errDo := http.DefaultClient.Do(req)

	// HTTP response errors and non-200 status codes will
	// 1) log to an error file, and
	// 2) return a nil buffer and the error from http.DefaultClient
	if errDo != nil {
		Error.Printf("URL '%s': response error: %v\n", get.String(), errDo)
		return nil, errDo
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("URL '%s': status code error: [%d] %s", get.String(), res.StatusCode, res.Status)
	}

	buf := new(bytes.Buffer)
	if _, errRead := buf.ReadFrom(res.Body); errRead != nil {
		return nil, errRead
	}

	return buf, nil
}
