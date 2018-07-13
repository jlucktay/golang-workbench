package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

func getResponse(get url.URL) io.Reader {
	req, reqErr := http.NewRequest("GET", get.String(), nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	req.Header.Add("User-Agent", "jlucktay (dotfiles)")
	req.SetBasicAuth("jlucktay", ghpaToken)

	res, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("URL '%s': status code error: %d %s", get.String(), res.StatusCode, res.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	return buf
}
